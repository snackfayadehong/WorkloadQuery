package main

import (
	"WorkloadQuery/conf"
	clientDb "WorkloadQuery/db"
	"WorkloadQuery/logger"
	"WorkloadQuery/middleware"
	"WorkloadQuery/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	r := gin.New()
	r.Use(Cors()) // 跨域
	// defer logFile.Close()
	// r.Use(gin.LoggerWithConfig(*logConfig))
	// r.Use(gin.Recovery())
	white := conf.Configs.IPWhite.IPWhiteList
	r.Use(IPWhiteList(white), logger.GinLogger, logger.GinRecovery(true))
	err := r.SetTrustedProxies([]string{"172.21.1.158"})
	if err != nil {
		panic(err)
	}
	r.GET("/debug/pprof/*any", gin.WrapH(http.DefaultServeMux))
	router := r.Group("/api")
	{
		router.POST("/getWorkload", middleware.CheckTime, service.GetWorkload)
		router.POST("/getNoAccountEntry", middleware.CheckTime, service.GetNoAccountEntry)
		router.POST("/getUnCheckBills", middleware.CheckTime, service.GetUnCheckBills)
		router.POST("/getNoDeliveredPurchaseSummary", middleware.CheckTime, service.GetNoDeliveredPurchaseSummary)
		v1 := router.Group("/v1")
		{
			v1.POST("/change_prod", middleware.CheckRequestProdInfo, service.ChangeProductInfoService)
		}
	}
	// pprof
	go func() {
		if err := http.ListenAndServe("localhost:3008", nil); err != nil {
			zap.L().Error("ERROR", zap.Error(err))
		}
	}()
	err = r.Run(fmt.Sprintf("%s:%s", conf.Configs.Server.IP, conf.Configs.Server.Port))
	if err != nil {
		zap.L().Error("ERROR", zap.Error(err))
		return
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
		// c.Writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5173")
		// c.Writer.Header().Set("Access-Control-Allow-Methods", "POST")
		// c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Content-Length")
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

// 初始化程序
func init() {
	// 读取配置文件
	err := conf.InitSetting()
	if err != nil {
		panic(err)
	}
	err = logger.InitLog()
	if err != nil {
		return
	}
	err = clientDb.Init()
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
