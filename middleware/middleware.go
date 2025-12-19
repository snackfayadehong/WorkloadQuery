package middleware

import (
	"WorkloadQuery/controller"
	clientDb "WorkloadQuery/db"
	http2 "WorkloadQuery/http"
	"WorkloadQuery/logger"
	"WorkloadQuery/model"
	"WorkloadQuery/service"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
		logMsg := fmt.Sprintf("\r\n事件:接口返回\r\n出参:%s\r\n%s\r\n", res.Message, logger.LoggerEndStr)
		logger.AsyncLog(logMsg)
		c.Abort()
	}
	for _, v := range *req.C {
		err = checkCode(v)
		if err != nil {
			res.Code = 1
			res.Message = err.Error()
			c.JSON(http.StatusCreated, res)
			logMsg := fmt.Sprintf("\r\n事件:接口返回\r\n出参:%s\r\n%s\r\n", res.Message, logger.LoggerEndStr)
			logger.AsyncLog(logMsg)
			c.Abort()
			break
		}
	}
}

// 校验104分类和院内代码
// 增加修改人員信息檢查 2024-07-26
func checkCode(element model.ChangeInfoElement) error {
	var row int
	var sql = "select count(1) from TB_Employee where HisCode = ?"
	if db := clientDb.DB.Raw(sql, element.HRCode).Find(&row); db.Error != nil {
		return db.Error
	}
	if row == 0 {
		return fmt.Errorf("人员信息错误")
	}
	if len(element.Code) != 14 {
		return fmt.Errorf("失败,不正确的院内代码:%s", element.Code)
	}
	if element.CategoryCode != "" {
		if len(element.CategoryCode) != 10 || !strings.HasSuffix(element.CategoryCode, "0000") {
			return fmt.Errorf("失败;院内代码:%s,104分类:%s非三级目录", element.Code, element.CategoryCode)
		}
	}
	return nil
}

// 数据校验
// 1. 校验104分类

// CheckTime 校验请求数据是否为合法时间
func CheckTime(c *gin.Context) {
	res := http2.NewBaseResponse()
	var qt QueryTime
	err := c.ShouldBindBodyWith(&qt, binding.JSON)
	if err != nil {
		res.Code = 1
		res.Message = err.Error()
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}
	v1, err := time.Parse("2006-01-02 15:04:05", qt.Start)
	v2, err2 := time.Parse("2006-01-02 15:04:05", qt.End)
	if err != nil || err2 != nil || v1.IsZero() || v2.IsZero() {
		res.Code = 1
		res.Message = "无效请求数据,请核查请求时间"
		c.JSON(http.StatusBadRequest, res)
		c.Abort()
		return
	}
	// 校验后的时间存入上下文
	c.Set("startTime", qt.Start)
	c.Set("endTime", qt.End)
	c.Next()
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
