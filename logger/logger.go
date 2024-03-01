package logger

import (
	"WorkloadQuery/conf"
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat/go-file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"net"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

const LoggerEndStr = "----------------------------------------------------------------------------"
const InfoLevel = zapcore.InfoLevel
const ErrorLevel = zapcore.ErrorLevel

// InitLog 日志
func InitLog() (err error) {
	_, err = os.Stat("logs")
	if os.IsNotExist(err) {
		err = os.Mkdir("logs", 0700)
	}
	// 创建Core三大件，进行初始化
	encoder := getEncoder()
	// 判断日志级别
	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { // error级别
		return lev >= zap.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { // info和debug级别,debug级别是最低的
		return lev < zap.ErrorLevel && lev >= zap.DebugLevel
	})
	infoSyncer := getLogWriter("logs/", InfoLevel)
	errorSyncer := getLogWriter("logs/", ErrorLevel)
	// 创建核心-->如果是debug模式，就在控制台和文件都打印，否则就只写到文件中
	var core zapcore.Core
	if conf.Configs.Server.RunModel == "debug" {
		// 开发模式，日志输出到终端
		// NewTee创建一个核心，将日志条目复制到两个或多个底层核心中。
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), InfoLevel),
		)
	} else {
		coreInfo := zapcore.NewCore(encoder, infoSyncer, lowPriority)
		coreErr := zapcore.NewCore(encoder, errorSyncer, highPriority)
		core = zapcore.NewTee(coreInfo, coreErr)
	}
	// 创建 logger 对象
	// Warn()/Error()等级别的日志会输出堆栈，Debug()/Info()这些级别不会
	log := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.WarnLevel))
	// 替换全局的 logger, 后续在其他包中只需使用zap.L()调用即可
	zap.ReplaceGlobals(log)
	return
}

// 获取Encoder，给初始化logger使用的
func getEncoder() zapcore.Encoder {
	// 使用zap提供的 NewProductionEncoderConfig
	encoderConfig := zap.NewProductionEncoderConfig()
	// 设置时间格式
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	// 时间的key
	encoderConfig.TimeKey = "time"
	// 级别
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// 显示调用者信息
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	// 返回json 格式的 日志编辑器
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// 获取切割的问题，给初始化logger使用的
func getLogWriter(filename string, leavel zapcore.Level) zapcore.WriteSyncer {
	var logFileName string
	switch leavel {
	case zapcore.ErrorLevel:
		logFileName = filename + "ERROR_%Y%m%d.log"
	case zapcore.InfoLevel:
		logFileName = filename + "%Y%m%d.log"
	default:
		logFileName = filename + "Other_%Y%m%d.log"
	}
	hook, _ := rotatelogs.New(
		logFileName,
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*30),
		rotatelogs.WithRotationTime(time.Hour*24))
	return zapcore.AddSync(hook)
}

// GinLogger 用于替换gin框架的Logger中间件，不传参数，直接这样写
func GinLogger(c *gin.Context) {
	sugar := zap.L().Sugar()
	start := time.Now()
	// 在处理请求前获取body
	// 2024-1-26 以前获取入参的方法Make 1024的buf会导致入参过长时 参数不完整
	// buf := make([]byte, 1024)
	// n, _ := c.Request.Body.Read(buf)
	// // 去除转义字符
	// reqBody := string(buf[0:n])
	// r := strings.NewReplacer(" ", "", "\r", "", "\n", "", "\"", "")
	// reqData := r.Replace(reqBody)
	// c.Request.Body = io.NopCloser(bytes.NewBuffer(buf)) // 将读取后的字节流重新放入body 避免后续程序取不到body参数
	// 方法2
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = io.ReadAll(c.Request.Body)
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	reqData := string(bodyBytes)
	// // 方法3
	// w := middleware.ResponseWriter{
	// 	ResponseWriter: c.Writer,
	// 	B:              bytes.NewBuffer([]byte{}),
	// }
	// c.Writer = w
	// 视图函数执行完成，统计时间，记录日志
	cost := time.Since(start)
	sugar.Infof("\r\n事件:接口调用\r\nIP：%s\r\nURL：%s\r\nMethod：%s\r\n入参：%s\r\nError：%s\r\nCost：%s\r\n%s\r\n",
		c.ClientIP(), c.Request.URL.Path, c.Request.Method, reqData,
		c.Errors.ByType(gin.ErrorTypePrivate).String(), cost, LoggerEndStr)
	c.Next() // 执行视图函数
}

// GinRecovery 用于替换gin框架的Recovery中间件，因为传入参数，再包一层
func GinRecovery(stack bool) gin.HandlerFunc {
	logger := zap.L()
	return func(c *gin.Context) {
		defer func() {
			// defer 延迟调用，出了异常，处理并恢复异常，记录日志
			if err := recover(); err != nil {
				//  这个不必须，检查是否存在断开的连接(broken pipe或者connection reset by peer)---------开始--------
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				// http util包预先准备好的DumpRequest方法
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// 如果连接已断开，我们无法向其写入状态
					c.Error(err.(error))
					c.Abort()
					return
				}
				//  这个不必须，检查是否存在断开的连接(broken pipe或者connection reset by peer)
				// 是否打印堆栈信息，使用的是debug.Stack()，传入false，在日志中就没有堆栈信息
				if stack {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				// 有错误，直接返回给前端错误，前端直接报错
				// c.AbortWithStatus(http.StatusInternalServerError)
				// 该方式前端不报错
				c.String(200, "访问出错了")
			}
		}()
		c.Next()
	}
}
