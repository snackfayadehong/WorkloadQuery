package main

import (
	clientDb "WorkloadQuery/db"
	"WorkloadQuery/logger"
	"WorkloadQuery/middleware"
	"WorkloadQuery/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	r := gin.New()
	r.Use(Cors()) // 跨域
	// defer logFile.Close()
	// r.Use(gin.LoggerWithConfig(*logConfig))
	// r.Use(gin.Recovery())
	r.Use(logger.GinLogger, logger.GinRecovery(true))
	router := r.Group("/api")
	{
		router.POST("/getWorkload", middleware.CheckTime, service.GetWorkload)
	}
	err := r.Run(fmt.Sprintf("%s:%s", clientDb.Configs.Server.IP, clientDb.Configs.Server.Port))
	if err != nil {
		zap.L().Error("ERROR", zap.Error(err))
		return
	}
}

// Cors 开启跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.Writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5173")
		// c.Writer.Header().Set("Access-Control-Allow-Methods", "POST")
		// c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Content-Length")
		// c.Header("Access-Control-Allow-Origin", "http://127.0.0.1:5173")
		c.Header("Access-Control-Allow-Origin", "http://172.21.1.158:5173")
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
	err := clientDb.Init()
	if err != nil {
		panic(err)
	}
	err = logger.InitLog()
	if err != nil {
		zap.L().Error("ERROR:", zap.Error(err))
		return
	}
}
