package service

import (
	"WorkloadQuery/controller"
	"WorkloadQuery/logger"
	"fmt"
	"time"
)

// DeliveryRetryService 领用出库单重试
func DeliveryRetryService() {
	var logMsg string
	var d controller.DeliveryRequestInfo
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

// WrappedTask 定时执行
func WrappedTask() {
	var logMsg string
	now := time.Now()
	hour := now.Hour()
	if hour >= 8 && hour < 18 {
		logMsg = fmt.Sprintf("\r\n事件：定时任务\r\n时间：%v\r\n%s\r\n", now.Format("2006-01-02 15:04:05"), logger.LoggerEndStr)
		logger.AsyncLog(logMsg)
		DeliveryRetryService()
		logMsg = fmt.Sprintf("\r\n事件：定时任务\r\n时间：%v\r\n执行结果：执行完毕\r\n%s\r\n",
			now.Format("2006-01-02 15:04:05"), logger.LoggerEndStr)
		logger.AsyncLog(logMsg)
	} else {
		logMsg = fmt.Sprintf("\r\n事件：定时任务\r\n时间：%v\r\n\r\n执行结果：非执行时间直接跳过\r\n%s\r\n",
			now.Format("2006-01-02 15:04:05"), logger.LoggerEndStr)
	}
}
