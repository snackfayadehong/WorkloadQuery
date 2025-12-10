package service

import (
	"WorkloadQuery/controller"
	"WorkloadQuery/logger"
	"fmt"
)

// DeliveryRetryService 领用出库单重试
func DeliveryRetryService() {
	var logMsg string
	d := controller.DeliveryRequestInfo{
		Count: new(int64),
	}
	err := d.GetDeliveryNo()
	if err != nil {
		return
	}
	if *d.Count == 0 {
		return
	}
	logMsg = fmt.Sprintf("\r\n事件:查询领用出库失败业务数据\r\n查询结果:%+v\r\n%v\r\n%s\r\n", *d.Count, *d.De, logger.LoggerEndStr)
	logger.AsyncLog(logMsg)
	//接口调用
	err = d.DeliveryNoRetryToHis()
	if err != nil {
		logMsg = fmt.Sprintf("\r\n:事件：接口调用错误:%v\r\n%s\r\n", err.Error(), logger.LoggerEndStr)
		return
	}
}
