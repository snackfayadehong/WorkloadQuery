package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"time"
)

type QueryTime struct {
	Start string `json:"startTime" binding:"required"`
	End   string `json:"endTime" binding:"required"`
}

// CheckTime 校验请求数据是否为合法时间
func CheckTime(c *gin.Context) {
	var qt QueryTime
	err := c.ShouldBindBodyWith(&qt, binding.JSON)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "无效请求数据",
			"Data": err,
		})
		c.Abort()
	}
	v1, err := time.Parse("2006-01-02", qt.Start)
	v2, err2 := time.Parse("2006-01-02", qt.End)
	if err == nil && err2 == nil {
		if v1.IsZero() || v2.IsZero() {
			c.JSON(http.StatusOK, gin.H{
				"msg":  "无效请求数据,请核查请求时间",
				"Data": "",
			})
			c.Abort()
		}
	}
}
