package middleware

import (
	"WorkloadQuery/controller"
	"WorkloadQuery/logger"
	"WorkloadQuery/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

type QueryTime struct {
	Start string `json:"startTime" binding:"required"`
	End   string `json:"endTime" binding:"required"`
}

func CheckRequestProdInfo(c *gin.Context) {
	var req controller.RequestInfo
	var res service.Response
	err := c.ShouldBindBodyWith(&req.C, binding.JSON)
	if err != nil {
		res.Code = 1
		res.Message = err.Error()
		c.JSON(http.StatusCreated, res)
		zap.L().Sugar().Infof("\r\n事件:接口返回\r\n出参:%s\r\n%s\r\n", res.Message, logger.LoggerEndStr)
		c.Abort()
	}
	for _, v := range req.C {
		err = checkCode(v)
		if err != nil {
			res.Code = 1
			res.Message = err.Error()
			c.JSON(http.StatusCreated, res)
			zap.L().Sugar().Infof("\r\n事件:接口返回\r\n出参:%s\r\n%s\r\n", res.Message, logger.LoggerEndStr)
			c.Abort()
			break
		}
	}
}

// 校验104分类和院内代码
func checkCode(element controller.ChangeInfoElement) error {
	if len(element.Code) != 14 {
		return fmt.Errorf("失败,不正确的院内代码:%s", element.Code)
	}
	if len(*element.CategoryCode) != 10 || !strings.HasSuffix(*element.CategoryCode, "0000") {
		return fmt.Errorf("失败;院内代码:%s,104分类:%s非三级目录", element.Code, *element.CategoryCode)
	}
	return nil
}

// 数据校验
// 1. 校验104分类

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

// // 获取 body 方法1
// type ResponseWriter struct {
// 	gin.ResponseWriter
// 	B *bytes.Buffer
// }
//
// func (w ResponseWriter) Write(b []byte) (int, error) {
// 	// 向一个bytes.buffer中写一份数据来为获取body使用
// 	w.B.Write(b)
// 	// 完成gin.Context.Writer.Write()原有功能
// 	return w.ResponseWriter.Write(b)
// }
