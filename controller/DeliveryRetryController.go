package controller

import (
	clientDb "WorkloadQuery/db"
	"WorkloadQuery/hisInterface"
	"WorkloadQuery/logger"
	"WorkloadQuery/model"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// DeliveryRequestInfo 接口入参
type DeliveryRequestInfo struct {
	Count *int64
	De    *[]model.DeliveryNo
}

// DeliveryResponseInfo 接口出参
type DeliveryResponseInfo struct {
	hisInterface.KLBRBaseResponse
	Data hisInterface.DeliveryData `json:"data"`
}

func (d *DeliveryRequestInfo) processSingleDelivery(raw model.DeliveryNo) error {
	// 准备请求数据
	raw.Ckdh += raw.DetailSort
	data, err := json.Marshal(raw)
	if err != nil {
		return fmt.Errorf("序列化请求数据失败: %w", err)
	}

	// 构建请求
	k := hisInterface.KLBRRequest{
		Headers: hisInterface.NewReqHeaders("herp-clckgl"),
		Url:     hisInterface.BaseUrl + "herp-clckgl/1.0",
		ReqData: data,
	}
	// 发送HTTP请求
	res, err := k.KLBRHttpPost()
	if err != nil {
		return fmt.Errorf("HTTP请求失败: %w", err)
	}

	// 解析响应
	var HisRes DeliveryResponseInfo
	if err = json.Unmarshal(*res, &HisRes); err != nil {
		logMsg := fmt.Sprintf("\r\n事件:接口请求跟踪\r\n出参：%v\r\n%s\r\n", res, logger.LoggerEndStr)
		logger.AsyncLog(logMsg)
		return fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查响应状态
	if HisRes.AckCode != "200.1" {
		return fmt.Errorf("接口返回错误1: %s", HisRes.AckMessage)
	}

	if len(HisRes.Data.Fhxx) == 0 {
		return fmt.Errorf("响应数据中缺少Fhxx信息")
	}

	fhxx := HisRes.Data.Fhxx[0]
	if strings.TrimSpace(fhxx.Ckdh) == "" && strings.TrimSpace(fhxx.Sczt) != "0" {
		return fmt.Errorf("接口返回错误2: %s", fhxx.Scsm)
	}

	// 记录成功日志
	logMsg := fmt.Sprintf("\r\n事件:接口请求跟踪\r\n出参:%+v\r\n%s\r\n", HisRes.Data, logger.LoggerEndStr)
	logger.AsyncLog(logMsg)

	// 更新数据库
	tx := clientDb.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	deliveryID := raw.Ckdh[:len(raw.Ckdh)-1]
	hisCkdh := fhxx.Ckdh

	if db := tx.Exec(clientDb.UpdateDelivery_Sql, hisCkdh, deliveryID, raw.DetailSort); db.Error != nil {
		tx.Rollback()
		return fmt.Errorf("更新数据库失败: %w", db.Error)
	}
	if err := tx.Commit(); err != nil {
		return err.Error
	}
	return nil
}

func (d *DeliveryRequestInfo) DeliveryNoRetryToHis() (err error) {
	var successCount int64
	for _, raw := range *d.De {
		if err = d.processSingleDelivery(raw); err != nil {
			// 记录错误但继续处理其他项目
			logMsg := fmt.Sprintf("\r\n事件:处理配送单失败\r\n配送单: %s\r\n错误: %v\r\n%s\r\n",
				raw.Ckdh, err, logger.LoggerEndStr)
			logger.AsyncLog(logMsg)
			continue
		}
		successCount++
	}

	if successCount == 0 {
		return fmt.Errorf("所有配送单处理失败")
	}

	logMsg := fmt.Sprintf("\r\n事件:配送单处理完成\r\n成功数量: %d\r\n总数量: %d\r\n%s\r\n",
		successCount, len(*d.De), logger.LoggerEndStr)
	logger.AsyncLog(logMsg)
	return nil
}

func (d *DeliveryRequestInfo) GetDeliveryNo() (err error) {
	var now = time.Now()
	s := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()) // 当天0时0点
	e := now.Add(-10 * time.Minute)                                                // 当前时间前推10分钟
	startDate := s.Format("2006-01-02 15:04:05")
	endDate := e.Format("2006-01-02 15:04:05")
	db := clientDb.DB.Raw(clientDb.QueryBillNo, startDate, endDate).Find(&d.De)
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected == 0 {
		*d.Count = 0
		logMsg := fmt.Sprintf("\r\n事件:查询领用出库失败业务数据\r\n查询结果:无数据记录\r\n%s\r\n", logger.LoggerEndStr)
		logger.AsyncLog(logMsg)
	} else {
		*d.Count = db.RowsAffected
	}
	return nil
}
