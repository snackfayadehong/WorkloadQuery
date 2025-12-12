package service

import (
	"WorkloadQuery/controller"
	"WorkloadQuery/logger"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Response struct {
	Code    int
	Message string
}

func ChangeProductInfoService(c *gin.Context) {
	var res Response
	// 查询条件
	var where []string
	// 入参
	var req controller.RequestInfo
	_ = c.ShouldBindBodyWith(&req.C, binding.JSON)
	// 将入参多条数据 Code 整合为一个 where 条件
	var seen = make(map[string]bool)
	for _, v := range *req.C {
		if !seen[v.Code] {
			seen[v.Code] = true
			where = append(where, v.Code)
		} else {
			res.Code = 1
			res.Message = fmt.Sprintf("%s入参存在多条;", v.Code)
			logMsg := fmt.Sprintf("\r\n事件:接口返回\r\n出参:%s\r\n%s\r\n", res.Message, logger.LoggerEndStr)
			logger.AsyncLog(logMsg)
			c.JSON(http.StatusCreated, res)
			return
		}
	}
	// 获取怡道系统产品基本信息
	prod, err := req.GetProductInfo(where)
	if err != nil {
		res.Code = 1
		res.Message = err.Error()
		c.JSON(http.StatusCreated, res)
		logMsg := fmt.Sprintf("\r\n事件:接口返回\r\n出参:%s\r\n%s\r\n", res.Message, logger.LoggerEndStr)
		logger.AsyncLog(logMsg)
		return
	}
	err2 := req.ChangeMisProductInfo(prod, c.ClientIP())
	if err2 != nil {
		res.Code = 1
		res.Message = err2.Error()
		c.JSON(http.StatusCreated, res)
		logMsg := fmt.Sprintf("\r\n事件:接口返回\r\n出参:%s\r\n%s\r\n", res.Message, logger.LoggerEndStr)
		logger.AsyncLog(logMsg)
		return
	}
	c.JSON(http.StatusCreated, res)
	logMsg := fmt.Sprintf("\r\n事件:接口返回\r\n出参:%s\r\n%s\r\n", res.Message, logger.LoggerEndStr)
	logger.AsyncLog(logMsg)
}
