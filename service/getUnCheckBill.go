package service

import (
	"WorkloadQuery/controller"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

func GetUnCheckBills(c *gin.Context) {
	UnCheckBillTimeInterval := controller.QueryUncheckTime{}
	_ = c.ShouldBindBodyWith(&UnCheckBillTimeInterval, binding.JSON)
	UnCheckBillTimeInterval.StartTime += " 00:00:00.000"
	UnCheckBillTimeInterval.EndTime += " 23:59:59.000"
	unCheckBills := UnCheckBillTimeInterval.UnCheckBillQuery()
	if unCheckBills == nil || len(*unCheckBills) == 0 {
		c.JSON(http.StatusNoContent, gin.H{
			"msg":  "无数据",
			"Data": "",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "成功",
			"Data": unCheckBills,
		})
	}
}
