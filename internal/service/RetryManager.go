package service

import (
	"SupperSystem/internal/controller"
	"SupperSystem/internal/model"
	http2 "SupperSystem/pkg/http"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type RetryListRequest struct {
	QueryType string `json:"queryType" binding:"required"`
	StartTime string `json:"startTime" binding:"required"`
	EndTime   string `json:"endTime" binding:"required"`
}
type RetryExecuteRequest struct {
	Type       string `json:"type" binding:"required"`
	BillNo     string `json:"billno" binding:"required"`
	DetailSort string `json:"detailSort"`
}

// HandleRetryList 获取待重试列表
func HandleRetryList(c *gin.Context) {
	var req RetryListRequest
	res := http2.NewBaseResponse()
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.QueryType == "delivery" {
		d := &controller.DeliveryRequestInfo{
			Count: new(int64),
			De:    &[]model.DeliveryNo{},
		}
		if err := d.GetDeliveryNo(req.StartTime, req.EndTime); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		res.Message = "查询成功"
		res.Data = d.De
		c.JSON(http.StatusOK, res)
	} else if req.QueryType == "refund" {
		r := &controller.RefundRequestInfo{
			Count: new(int64),
			Re:    &[]model.RefundNo{},
		}
		if err := r.GetRefundNo(req.StartTime, req.EndTime); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		res.Message = "查询成功"
		res.Data = r.Re
		c.JSON(http.StatusOK, res)
	}
}

// HandleRetryExecute 重试
func HandleRetryExecute(c *gin.Context) {
	var req RetryExecuteRequest
	res := http2.NewBaseResponse()
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 方式为出库
	if req.Type == "delivery" {
		d := controller.DeliveryRequestInfo{
			Count: new(int64),
			De: &[]model.DeliveryNo{
				{
					Ckdh:       req.BillNo,
					DetailSort: req.DetailSort,
					Ckfs:       "01",
				},
			},
		}
		// 传HIS
		if err := d.DeliveryNoRetryToHis(); err != nil {
			res.Message = err.Error()
			res.Code = 1
			c.JSON(http.StatusInternalServerError, res)
			return
		}
		res.Message = "处理成功"
		c.JSON(http.StatusOK, res)
	} else if req.Type == "refund" { // 方式为科室退库
		r := controller.RefundRequestInfo{
			Count: new(int64),
			Re: &[]model.RefundNo{
				{
					Yddh: req.BillNo,
					Rkfs: "02",
				},
			},
		}
		if err := r.RetryRefundToHis(); err != nil { // 传HIS
			res.Message = err.Error()
			res.Code = 1
			c.JSON(http.StatusInternalServerError, res)
			return
		}
		res.Message = "处理成功"
		c.JSON(http.StatusOK, res)
	}
}
