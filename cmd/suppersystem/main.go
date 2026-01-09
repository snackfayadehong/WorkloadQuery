package main

import (
	"SupperSystem/configs"
	task "SupperSystem/internal/Task"
	"SupperSystem/internal/service"
	clientDb "SupperSystem/pkg/db"
	"SupperSystem/pkg/logger"
	"SupperSystem/pkg/middleware"
	"context"
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {

	// 初始化基础组件
	initApplication()
	// 初始化TaskManager
	taskManager := task.NewTaskManager()
	// 启动定时任务
	if err := taskManager.Start(); err != nil {
		zap.L().Fatal("定时任务启动失败", zap.Error(err))
	}
	// 程序退出清理资源
	defer cleanupResources(taskManager)
	// 初始化 Gin引擎和路由
	r := setupGinEngine()
	setupRoutes(r)
	// 启动辅助服务
	//startAuxiliaryServices(r)
	// 主服务
	startMainService(r)
}

// 初始化应用程序2
func initApplication() {
	runtime.SetBlockProfileRate(1)
	// 读取配置文件
	rootFile, _ := exec.LookPath(os.Args[0])
	path, err := filepath.Abs(rootFile)
	if err != nil {
		panic(fmt.Sprintf("获取程序路径失败:%v", err))
	}
	index := strings.LastIndex(path, string(os.PathSeparator))
	rootPath := path[:index]
	// 初始化配置
	err = conf.InitSetting(rootPath)
	if err != nil {
		panic(fmt.Sprintf("初始化配置失败:%v", err))
	}
	// 初始化日志
	err = logger.InitLog()
	if err != nil {
		panic(fmt.Sprintf("初始化日志失败:%v", err))
	}
	// 初始化数据库
	err = clientDb.Init(rootPath)
	if err != nil {
		panic(err)
	}
	// 根据配置文件设置选择程序环境
	switch conf.Configs.Server.RunModel {
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}

// 清理资源
func cleanupResources(taskManager *task.TaskManager) {
	// 停止定时
	taskManager.Stop()
	// 关闭日志
	logger.Close()
}

// Gin引擎
func setupGinEngine() *gin.Engine {
	r := gin.New()
	// 信任代理
	err := r.SetTrustedProxies([]string{"172.21.1.158"})
	if err != nil {
		panic(err)
	}
	// 中间件
	white := conf.Configs.IPWhite.IPWhiteList
	r.Use(Cors())
	r.Use(IPWhiteList(white), logger.GinLogger, logger.GinRecovery(true))
	return r
}

// 设置路由
func setupRoutes(r *gin.Engine) {
	// pprof
	r.GET("/debug/pprof/*any", gin.WrapH(http.DefaultServeMux))
	// 健康检查路由
	r.GET("/health", func(c *gin.Context) {
		c.JSONP(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
			"service":   "SupperSystem",
		})
	})

	// API 路由组
	router := r.Group("/api")
	{
		router.POST("/getWorkload", middleware.CheckTime, service.WorkloadServiceInstance.HandleWorkloadRequest)
		router.POST("/dict/compare", service.DictCompareServiceInstance.HandleCompareRequest)
		v1 := router.Group("/v1")
		{
			v1.POST("/change_prod", middleware.CheckRequestProdInfo, service.ChangeProductInfoService)
		}
		retry := router.Group("/retry")
		{
			retry.POST("/list", middleware.CheckTime, service.HandleRetryList) // 获取失败列表
			retry.POST("/execute", service.HandleRetryExecute)                 // 执行重试动作
		}
	}
}

// 主服务
func startMainService(r *gin.Engine) {
	// create Http
	addr := fmt.Sprintf("%s:%s", conf.Configs.Server.IP, conf.Configs.Server.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// start goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	// 关机
	setupGracefulShutdown(server)
}

// 关机
func setupGracefulShutdown(server *http.Server) {
	// 等待终端信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 阻塞直到收到信号
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// 关机
	if err := server.Shutdown(ctx); err != nil {
		zap.L().Error("Http服务关闭失败", zap.Error(err))
	}
}

// IPWhiteList Ip白名单
func IPWhiteList(w []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求的 IP 地址 无代理返回客户端IP,有代理返回代理IP
		ip := c.RemoteIP()
		// 检查 IP 地址是否在白名单中
		allowed := false
		for _, value := range w {
			if value == ip {
				allowed = true
				break
			}
		}
		// 如果 IP 地址不在白名单中，则返回错误信息
		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "IP address not allowed"})
			return
		}
		// 允许请求继续访问后续的处理函数
		c.Next()
	}
}

// Cors 开启跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		if conf.Configs.Server.RunModel == "debug" {
			c.Header("Access-Control-Allow-Origin", "http://127.0.0.1:5173")
		} else {
			c.Header("Access-Control-Allow-Origin", "http://172.21.1.158")
		}
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
		} else {
			c.Next()
		}
	}
}
