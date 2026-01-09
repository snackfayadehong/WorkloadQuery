package controller

import (
	"SupperSystem/internal/model"
	clientDb "SupperSystem/pkg/db"
	"SupperSystem/pkg/integration"
	"SupperSystem/pkg/logger"
	"encoding/json"
	"fmt"
	"strings"
)

type RefundRequestInfo struct {
	Count *int64
	Re    *[]model.RefundNo
}

type RefundResponseInfo struct {
	integration.KLBRBaseResponse
	Data integration.RefundData
}

func (r *RefundRequestInfo) processSingleRefund(raw model.RefundNo) error {
	full := &model.RefundFullSerializer{RefundNo: &raw}
	raw2 := full.RefundSerialize()
	// 准备请求数据
	data, err := json.Marshal(raw2)
	if err != nil {
		return fmt.Errorf("序列化请求数据失败: %w", err)
	}

	// 构建请求
	k := integration.KLBRRequest{
		Headers: integration.NewReqHeaders("herp-clrkgl"),
		Url:     integration.BaseUrl + "herp-clrkgl/1.0",
		ReqData: data,
	}
	// 发送  HTTP 请求
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
	//if db := tx.Exec(clientDb.UpdateRefund_Sql, raw.Yddh); db.Error != nil {
	//	tx.Rollback()
	//	return fmt.Errorf("更新数据库失败: %w", db.Error)
	//}
	if err := tx.Table("TB_Refund").Where("RetWarhouCode = ?", raw.Yddh).Update("SendStatus", 1).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新数据失败: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return err.Error
	}
	return nil
}

func (r *RefundRequestInfo) RetryRefundToHis() (err error) {
	var successCount int64
	for _, raw := range *r.Re {
		if err = r.processSingleRefund(raw); err != nil {
			// 记录错误但继续处理其他项目
			logMsg := fmt.Sprintf("\r\n事件:处理科室退库单失败\r\n退库单号: %s\r\n错误: %v\r\n%s\r\n",
				raw.Yddh, err, logger.LoggerEndStr)
			logger.AsyncLog(logMsg)
			continue
		}
		successCount++
	}

	if successCount == 0 {
		return fmt.Errorf("所有科室退库单处理失败")
	}

	logMsg := fmt.Sprintf("\r\n事件:科室退库单处理完成\r\n成功数量: %d\r\n总数量: %d\r\n%s\r\n",
		successCount, len(*r.Re), logger.LoggerEndStr)
	logger.AsyncLog(logMsg)
	return nil
}
func (r *RefundRequestInfo) GetRefundNo(startDate, endDate string) (err error) {
	//db := clientDb.DB.Raw(clientDb.QueryRefundBillno, startDate, endDate).Find(&r.Re)
	db := clientDb.DB.Table("TB_Refund a").
		Select("a.RetWarhouCode as yddh,'02' as rkfs,store.DepartmentName as storeHouseName,dept.DepartmentName as leaderDepartName").
		Joins("Left Join TB_Department store on store.DeptCode = a.TargetStorehouseID").
		Joins("Left Join TB_Department dept on dept.DeptCode = a.DeptCode").
		Where("ISNULL(a.SendStatus,?) = ?", "", "").
		Where("a.Status = ?", 51).
		Where("a.CreateTime >= ? And a.CreateTime < ?", startDate, endDate).
		Find(&r.Re)
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected == 0 {
		*r.Count = 0
		logMsg := fmt.Sprintf("\r\n事件:查询科室退库失败业务数据\r\n查询结果:无数据记录\r\n%s\r\n", logger.LoggerEndStr)
		logger.AsyncLog(logMsg)
	} else {
		*r.Count = db.RowsAffected
	}
	return nil
}
