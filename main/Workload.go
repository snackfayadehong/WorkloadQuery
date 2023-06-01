package main

import (
	clientDb "WorkloadQuery/db"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Employee struct {
	HRCode       int
	EmployeeName string
}

func main() {
	user := Employee{}
	r := gin.Default()
	r.Use(Cors()) // 跨域
	r.POST("/api/post", func(c *gin.Context) {
		clientDb.ProdAcceptQuery("2023-04-01", "2023-05-01")
		clientDb.DB.Select("HRCode,EmployeeName").First(&user)
		u, _ := json.Marshal(user)
		c.String(http.StatusOK, string(u))
	})
	err := r.Run("127.0.0.1:3007")
	if err != nil {
		return
	}
}

// Cors 开启跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.Writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5173")
		// c.Writer.Header().Set("Access-Control-Allow-Methods", "POST")
		// c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Content-Length")
		c.Header("Access-Control-Allow-Origin", "http://127.0.0.1:5173")
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

// 数据库连接
func init() {
	err := clientDb.InitDb()
	if err != nil {
		return
	}
}
