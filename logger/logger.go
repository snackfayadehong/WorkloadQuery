package logger

import (
	clientDb "WorkloadQuery/db"
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

const level = zapcore.InfoLevel

// InitLog 日志
func InitLog() (logFile *os.File, logConfig *gin.LoggerConfig, err error) {
	_, err = os.Stat("logs")
	if os.IsNotExist(err) {
		err = os.Mkdir("logs", 0700)
	}
	// 创建Core三大件，进行初始化
	writeSyncer := getLogWriter("logs/")
	encoder := getEncoder()
	// 创建核心-->如果是debug模式，就在控制台和文件都打印，否则就只写到文件中
	var core zapcore.Core
	if clientDb.Configs.Server.RunModel == "debug" {
		// 开发模式，日志输出到终端
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		// NewTee创建一个核心，将日志条目复制到两个或多个底层核心中。
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writeSyncer, level),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), level),
		)
	} else {
		core = zapcore.NewCore(encoder, writeSyncer, level)
	}

	// core := zapcore.NewCore(encoder, writeSyncer, level)
	// 创建 logger 对象
	// Warn()/Error()等级别的日志会输出堆栈，Debug()/Info()这些级别不会
	log := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.WarnLevel))
	// 替换全局的 logger, 后续在其他包中只需使用zap.L()调用即可
	zap.ReplaceGlobals(log)
	return
	// logTime := time.Now().Format("2006-01-02")
	// _, err = os.Stat(fmt.Sprintf("logs/%s.log", logTime))
	// if os.IsNotExist(err) {
	// 	logFile, err = os.OpenFile(fmt.Sprintf("logs/%s.log", logTime), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	// 	if err != nil {
	// 		return
	// 	}
	// }
	// // gin日志配置
	// logConfig = &gin.LoggerConfig{
	// 	Formatter: func(params gin.LogFormatterParams) string {
	// 		return fmt.Sprintf("客户端IP:%s,请求时间:[%s],请求方式:%s,请求地址:%s,http协议版本:%s,请求状态码:%d,响应时间:%s,客户端:%s,错误信息:%s\r\n",
	// 			params.ClientIP,
	// 			params.TimeStamp.Format("2006-01-02 15:04:05"),
	// 			params.Method,
	// 			params.Path,
	// 			params.Request.Proto,
	// 			params.StatusCode,
	// 			params.Latency,
	// 			params.Request.UserAgent(),
	// 			params.ErrorMessage,
	// 		)
	// 	},
	// 	Output: logFile,
	// }
}

// 获取Encoder，给初始化logger使用的
func getEncoder() zapcore.Encoder {
	// 使用zap提供的 NewProductionEncoderConfig
	encoderConfig := zap.NewProductionEncoderConfig()
	// 设置时间格式
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	// 时间的key
	encoderConfig.TimeKey = "time"
	// 级别
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// 显示调用者信息
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	// 返回json 格式的 日志编辑器
	return zapcore.NewJSONEncoder(encoderConfig)
}

// 获取切割的问题，给初始化logger使用的
func getLogWriter(filename string) zapcore.WriteSyncer {
	hook, _ := rotatelogs.New(
		filename+"%Y%m%d"+".log",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*30),
		rotatelogs.WithRotationTime(time.Hour*24))
	return zapcore.AddSync(hook)
}

// GinLogger 用于替换gin框架的Logger中间件，不传参数，直接这样写
func GinLogger(c *gin.Context) {
	logger := zap.L()
	start := time.Now()
	path := c.Request.URL.Path
	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)
	// 去除转义字符
	reqBody := string(buf[0:n])
	r := strings.NewReplacer(" ", "", "\r", "", "\n", "", "\"", "")
	reqData := r.Replace(reqBody)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(buf)) // 将读取后的字节流重新放入body 避免后续程序取不到body参数
	c.Next()                                            // 执行视图函数
	// 视图函数执行完成，统计时间，记录日志
	cost := time.Since(start)
	logger.Info(path,
		zap.Int("status", c.Writer.Status()),
		zap.String("method", c.Request.Method),
		zap.String("path", path),
		zap.String("query", reqData),
		zap.String("ip", c.ClientIP()),
		zap.String("user-agent", c.Request.UserAgent()),
		zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
		zap.Duration("cost", cost),
	)
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
