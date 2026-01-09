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

// DeliveryRequestInfo 接口入参
type DeliveryRequestInfo struct {
	Count *int64
	De    *[]model.DeliveryNo
}

// DeliveryResponseInfo 接口出参
type DeliveryResponseInfo struct {
	integration.KLBRBaseResponse
	Data integration.DeliveryData `json:"data"`
}

func (d *DeliveryRequestInfo) processSingleDelivery(raw model.DeliveryNo) error {
	full := &model.DeliveryFullSerializer{DeliveryNo: &raw} // 构造用于请求 HIS 的接口请求参数
	raw2 := full.DeliverySerialize()
	data, err := json.Marshal(raw2)
	if err != nil {
		return fmt.Errorf("序列化请求数据失败: %w", err)
	}

	// 构建请求
	k := integration.KLBRRequest{
		Headers: integration.NewReqHeaders("herp-clckgl"),
		Url:     integration.BaseUrl + "herp-clckgl/1.0",
		ReqData: data,
	}
	// 发送 HTTP请求
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
	hisCkdh := fhxx.Ckdh

	//if db := tx.Exec(clientDb.UpdateDelivery_Sql, hisCkdh, deliveryid, raw.DetailSort); db.Error != nil {
	//	tx.Rollback()
	//	return fmt.Errorf("更新数据库失败: %w", db.Error)
	//}
	if err := tx.Table("TB_DeliveryApplyDetailRecord").
		Where("DeliveryId = ? AND DetailSort = ?", raw.Ckdh, raw.DetailSort).
		Update("OutNumber", hisCkdh).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新数据库失败:%w", err)
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

func (d *DeliveryRequestInfo) GetDeliveryNo(startDate, endDate string) (err error) {
	db := clientDb.DB.Table("TB_DeliveryApplyDetailRecord As dr").
		Select("'01' AS ckfs, d.DeliveryID AS ckdh, dr.DetailSort AS detailSort,"+
			"d.DeliveryCode,DE.DepartmentName as leaderDepartName,DEPT.DepartmentName as storeHouseName").
		Joins("JOIN TB_DeliveryApply d on dr.DeliveryID = d.DeliveryID").
		Joins("Left Join TB_Department dept on dept.DeptCode = d.StorehouseID").
		Joins("Left Join TB_Department de on de.DeptCode = d.DeptCode").
		Where("dr.IsVoid = ?", 0).
		Where("d.Source = ?", "1").
		Where("d.IsStockGoods = ?", "0").
		Where("d.Type = ?", "1").
		Where("d.Status IN ?", []int{61, 71, 41, 81, 22, 91, 19, 29, 99}).
		Where("(d.IsStockGoods <> '1' OR d.IsStockGoods IS NULL)").
		Where("ISNULL(dr.OutNumber, '') = ?", "").
		Where("dr.UpdateTime >= ? AND dr.UpdateTime <= ?", startDate, endDate).
		Group("dr.DetailSort, d.DeliveryID,d.DeliveryCode,DEPT.DepartmentName,de.DepartmentName").
		Find(&d.De)
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
