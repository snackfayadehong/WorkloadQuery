package service

import (
	"WorkloadQuery/controller"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

func ChangeProductInfoService(c *gin.Context) {
	// 返回信息
	var code int
	var msg string
	// 查询条件
	var where string
	// 入参
	var rep []controller.ChangeInfoElement
	err := c.ShouldBindBodyWith(&rep, binding.JSON)
	if err != nil {
		panic(err)
	}
	fmt.Println(&rep)
	// 将入参多条数据Code整合为一个where条件
	var seen = make(map[string]bool)
	for _, v := range rep {
		if !seen[v.Code] {
			seen[v.Code] = true
			where += fmt.Sprintf("'%s',", v.Code)
		} else {
			code = 1
			msg += fmt.Sprintf("%s入参存在多条;", v.Code)
		}
	}
	if code == 1 {
		c.JSON(http.StatusOK, gin.H{
			"Code":    code,
			"Message": msg,
			"err":     err,
		})
		return
	} else {
		c.JSON(201, gin.H{
			"Code":    code,
			"Message": msg,
			"err":     err,
		})
	}
	fmt.Println("123")
	// // 删除最后一个 ','
	// where = where[:len(where)-1]
}
