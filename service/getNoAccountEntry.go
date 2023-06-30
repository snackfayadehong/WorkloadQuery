package service

import (
	"WorkloadQuery/controller"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

type Time struct {
	Start string `json:"startTime" binding:"required"`
	End   string `json:"endTime" binding:"required"`
}

func GetNoAccountEntry(c *gin.Context) {
	times := Time{}
	_ = c.ShouldBindBodyWith(&times, binding.JSON)
	times.End += " 23:59:59.000"
	_ = c.ShouldBindBodyWith(&times, binding.JSON)
	Ae := controller.NoAccountEntryQuery(&times.End)

	if Ae == nil || len(*Ae) == 0 {
		c.JSON(http.StatusNoContent, gin.H{
			"msg":  "无数据",
			"data": "",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "成功",
			"data": Ae,
		})
	}
}
