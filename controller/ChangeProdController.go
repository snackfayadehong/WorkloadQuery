package controller

import (
	clientDb "WorkloadQuery/db"
	"WorkloadQuery/hisInterface"
	"WorkloadQuery/logger"
	"WorkloadQuery/model"
	"WorkloadQuery/utity"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

/*
	根据老物资采购系统传过来的字典信息变更
*/

type RequestInfo struct {
	C *[]model.ChangeInfoElement
}

type klbrRes struct {
	hisInterface.KLBRBaseResponse
	Data hisInterface.ProductChangeData `json:"data"`
}

const UpdateCateCodeSql = "Update TB_ProductInfo Set CusCategoryCode = ? where ProductInfoID = ?"
const UpdateHospitalSpecSql = "Update TB_ProductInfo set HospitalSpec = ? where ProductInfoID = ?"
const UpdateHospitalNameSql = "Update TB_ProductInfo Set HisProductCode3 =? where ProductInfoID = ?"
const UpdateOpenTenderSql = "Update TB_ProductInfo set OpenTender = ? where ProductInfoID = ?"

func ChangeHisProductInfo(p model.ChangeInfoElement) error {
	var klbrres klbrRes
	// 1. 查询HIS需要的产品基本信息
	var his model.ChangeHisProductInfoModel
	if db := clientDb.DB.Raw(clientDb.ProductInfo_UpdatePostDataSQL, p.Code).Find(&his); db.Error != nil {
		return db.Error
	}
	if p.OpenTender == "1" {
		if !strings.HasPrefix(his.Ypmc, "(g)") {
			his.Ypmc = "(g)" + his.Ypmc
		}
		his.Jcsfzb = "1"
	} else if p.OpenTender == "0" {
		his.Jcsfzb = "0"
	}
	his.Kfbz = p.EighteenProdType // 18类重点监控耗材序号
	his.Xgczy = p.HRCode          // 修改人员工号
	//val := reflect.ValueOf(&his)
	//if val.Kind() == reflect.Ptr { // 检查 val 是否为指针
	//	val = val.Elem() // 获取指针指向的实际值
	//}
	//if val.Kind() == reflect.Struct { // 确保 val 是一个结构体
	//	fieldNum := val.NumField() // 获取该结构体有几个字段
	//	// 遍历结构体字段
	//	for i := 0; i < fieldNum; i++ {
	//		field := val.Field(i)
	//		if !field.CanSet() { // 检查字段是否可以设置
	//			return fmt.Errorf("%s", "value can't be set")
	//		}
	//		switch field.Kind() {
	//		case reflect.String:
	//			if field.String() == "" {
	//				field.SetString("20240815")
	//			}
	//		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	//			if field.Int() == 0 {
	//				field.SetInt(20240815)
	//			}
	//		default:
	//			continue
	//		}
	//	}
	//
	//}
	// 去除制表符
	utity.RemoveTabsFromStruct(&his)
	// 2. json序列化
	data, err := json.Marshal(his)
	if err != nil {
		return err
	}
	// 3. 接口调用
	k := hisInterface.KLBRRequest{}
	k.Headers = hisInterface.NewReqHeaders("MIS-proc-Edit-cljjxx")
	k.Url = hisInterface.BaseUrl + "MIS-proc-Edit-cljjxx/1.0"
	k.ReqData = data
	res, err := k.KLBRHttpPost()
	if err != nil {
		return err
	}
	// 4. 接口返回
	if err = json.Unmarshal(*res, &klbrres); err != nil {
		return err
	}
	//baseRep, fhxxList, err := hisInterface.ParseResPonse[map[string]interface{}](*res)
	//if err != nil {
	//	return fmt.Errorf("ParseResPonse err: %v", err)
	//}
	if klbrres.AckCode != "200.1" {
		return fmt.Errorf("%s", klbrres.AckMessage)
	} else {
		logMsg := fmt.Sprintf("\r\n事件:His接口返回\r\n出参:%s\r\n%s\r\n", klbrres.Data, logger.LoggerEndStr)
		logger.AsyncLog(logMsg)
	}
	return nil
}

/*
ChangeMisProductInfo
更改Mis产品基本信息
*/
func (i *RequestInfo) ChangeMisProductInfo(prod []model.ProductInfo, ip string) error {
	// 开启事务
	tx := clientDb.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// 查找入参和物资系统返回查询结果中相同的记录 时间复杂度O(M+N)
	// 使用Map存储prod切片中的Code
	pMap := make(map[string]int, len(prod))
	for index, v := range prod {
		pMap[v.Code] = index
	}
	// 在入参中找到Code相同的
	for _, item := range *i.C {
		if pIndex, ok := pMap[item.Code]; ok {
			var updateMsg string
			updateMsg = fmt.Sprintf("产品ID:%v", prod[pIndex].ProductInfoID)
			// 停供信息（1.直接修改科室领用白名单表状态，2.需要做停供处理的不处理字典信息)
			switch item.SupplyStatus {
			case "1":
				if err := UpdateProductSupplyStatus(tx, prod[pIndex], &updateMsg); err != nil {
					return err
				}
			case "0":
				// 修改字典信息的业务逻辑
				// 1. 104分类；
				if err := UpdateCategoryCode(tx, &item, prod[pIndex], &updateMsg); err != nil {
					return err
				}
				// 2. 更新院内产品名称、规格信息
				if err := UpdateHospitalInfo(tx, &item, prod[pIndex], &updateMsg); err != nil {
					return err
				}
				// 3. 判断集采审核状态并更新集采信息
				if err := UpdateYgcgidInfo(tx, &item, prod[pIndex], &updateMsg); err != nil {
					return err
				}
				// 4. 更新TradeCode流水号
				if err := UpdateTradeCodeInfo(tx, &item, prod[pIndex], &updateMsg); err != nil {
					return err
				}
				// 5. 更新医保代码
				if err := UpdateMedicareCodeInfo2(tx, &item, prod[pIndex], &updateMsg); err != nil {
					return err
				}
				// 6. 写入系统编码系统编号
				if err := UpdateJCSysInfo(tx, &item, prod[pIndex], &updateMsg); err != nil {
					return err
				}
				// 7. 更新集采状态
				if err := UpdateProductOpenTender(tx, &item, prod[pIndex], &updateMsg); err != nil {
					return err
				}
			default:
				return fmt.Errorf("入参信息有误")
			}
			// 写入日志表
			if updateMsg != "" && updateMsg != fmt.Sprintf("产品ID:%v", prod[pIndex].ProductInfoID) {
				if db := tx.Exec("Insert Into TB_PoductInfoChangeLog(Prod_Id,Context,IP,UpdateTime) values(?,?,?,?)",
					prod[pIndex].ProductInfoID, updateMsg, ip, time.Now()); db.Error != nil {
					tx.Rollback()
					return db.Error
				}
			}
			tx.Commit()
			// 请求HIS
			// 停供则处理下一个，不与HIS交互
			if item.SupplyStatus == "1" {
				continue
			}
			if err := ChangeHisProductInfo(item); err != nil {
				return err
			}
		}
	}
	return nil
}

/*
GetProductInfo
获取物资产品字典信息,返回不重复的字典信息
1. 相同院内编码多条记录，但供货状态PurState或isVoid不同的则跳过,只处理正常记录
2. 相同院内编码多条记录,供货状态相同,视为异常记录,接口报错返回
*/
func (i *RequestInfo) GetProductInfo(Where []string) ([]model.ProductInfo, error) {
	var prod []model.ProductInfo         // 原始记录
	var NoRepeatProd []model.ProductInfo // 保留不重复的记录
	var exception []model.ExceptionProd  // 异常信息记录
	var msg string
	db := clientDb.DB.Raw(clientDb.QueryProd, Where).Find(&prod)
	if db.Error != nil {
		return prod, db.Error
	}
	// 检查 查询结果中同一院内编码是否存在多条记录,且非停用或停供产品
	seen := make(map[string]bool)
	exceptionMap := make(map[string]bool)
	for _, el := range prod {
		// 只添加供货状态正常的记录
		if !seen[el.Code] && el.PurState == 0 && el.IsVoid == 0 {
			seen[el.Code] = true
			NoRepeatProd = append(NoRepeatProd, el)
			continue
		}
		// 记录异常信息
		exception = append(exception, model.ExceptionProd{Code: el.Code, PurState: el.PurState, IsVoid: el.IsVoid, HasError: true})
	}
	// 循环异常信息，当异常信息Code在seen正常中存在时，判断供货状态，如果非正常供货状态则跳过，否者说明此记录重复,记录并返回
	// 当异常信息在seen中不存在,则说明此记录供货关系异常,记录并返回
	for _, v := range exception {
		if seen[v.Code] && v.HasError {
			if v.PurState != 0 || v.IsVoid != 0 {
				continue
			}
		}
		// 处理异常信息
		if !exceptionMap[v.Code] {
			exceptionMap[v.Code] = true
			msg += fmt.Sprintf("%s有重复记录或供货关系异常", v.Code)
		}
	}
	if msg == "" {
		return NoRepeatProd, nil
	}
	return NoRepeatProd, fmt.Errorf(msg)
}

// UpdateProductOpenTender 更新产品集采状态
func UpdateProductOpenTender(tx *gorm.DB, item *model.ChangeInfoElement, prod model.ProductInfo, context *string) error {
	if item.OpenTender != "" && item.OpenTender != prod.OpenTender {
		if db := tx.Exec(UpdateOpenTenderSql, item.OpenTender, prod.ProductInfoID); db.Error != nil {
			tx.Rollback()
			return db.Error
		}
		*context += fmt.Sprintf("OpenTender:集采状态(%s)变更为(%s);",
			prod.OpenTender, item.OpenTender)
	}
	return nil
}

// UpdateCategoryCode 更新104分类
func UpdateCategoryCode(tx *gorm.DB, item *model.ChangeInfoElement, prod model.ProductInfo, context *string) error {
	// 入参中104分类为第3级,物资系统为第4级,查询物资系统第4级代码对应的第3级与入参是否相同
	// 如果不相同则修改物资为第3级,相同则不更新
	if item.CategoryCode != "" && item.CategoryCode != prod.ParentCusCategoryCode && item.CategoryCode != prod.CusCategoryCode {
		// 修改为第3级
		if db := tx.Exec(UpdateCateCodeSql, item.CategoryCode, prod.ProductInfoID); db.Error != nil {
			tx.Rollback()
			return db.Error
		}
		*context = fmt.Sprintf("CusCategoryCode(%s)变更为(%s);", prod.CusCategoryCode, item.CategoryCode)
	}
	return nil
}

// UpdateHospitalInfo 更新院内名称、院内规格
func UpdateHospitalInfo(tx *gorm.DB, item *model.ChangeInfoElement, prod model.ProductInfo, context *string) error {
	if item.HospitalSpec != "" && item.HospitalSpec != prod.HospitalSpec {
		if db := tx.Exec(UpdateHospitalSpecSql, item.HospitalSpec, prod.ProductInfoID); db.Error != nil {
			tx.Rollback()
			return db.Error
		}
		*context += fmt.Sprintf("HospitalSpec:院内规格(%s)变更为(%s);",
			prod.HospitalSpec, item.HospitalSpec)
	}
	if item.HospitalName != "" && item.HospitalName != prod.HospitalName {
		if db := tx.Exec(UpdateHospitalNameSql, item.HospitalName, prod.ProductInfoID); db.Error != nil {
			tx.Rollback()
			return db.Error
		}
		*context += fmt.Sprintf("HospitalName:院内名称(%s)变更为(%s);",
			prod.HospitalName, item.HospitalName)
	}
	return nil
}

// UpdateYgcgidInfo 更新网采平台产品ID
func UpdateYgcgidInfo(tx *gorm.DB, item *model.ChangeInfoElement, prod model.ProductInfo, context *string) error {
	// 1 已审核  null '' 0 为未审核
	if item.YGCGID != "" {
		// 2024-09-30 只要入参与HisProductCode7Source和HisProductCode7不相同 直接更新
		if prod.YGCGID != item.YGCGID || prod.HisProductCode7Source != item.YGCGID {
			sqlStr := `Update TB_ProductInfo Set HisProductCode7 = ?,HisProductCode7Source = ? Where ProductInfoId =?`
			sqlStr2 := `Update TB_ProductInfo Set HisProductCode7 = ?, HisProductCode7Source = ? ,HisProductCode7Status = '1' where ProductInfoID= ?`
			if prod.HisProductCode7Status == "1" {
				if db := tx.Exec(sqlStr, item.YGCGID, item.YGCGID, prod.ProductInfoID); db.Error != nil {
					tx.Rollback()
					return db.Error
				}
				*context += fmt.Sprintf("集采产品ID:HisProductCode7(%s),HisProductCode7Source(%s)变更为(%s)", prod.YGCGID, prod.HisProductCode7Source, item.YGCGID)
			} else {
				if db := tx.Exec(sqlStr2, item.YGCGID, item.YGCGID, prod.ProductInfoID); db.Error != nil {
					tx.Rollback()
					return db.Error
				}
				*context += fmt.Sprintf("集采产品ID:HisProductCode7(%s),HisProductCode7Source(%s)变更为(%s)并已审核", prod.YGCGID, prod.HisProductCode7Source, item.YGCGID)
			}
		}
		//	if prod.HisProductCode7Status == "1" && item.YGCGID != prod.YGCGID {
		//		if db := tx.Exec("Update TB_ProductInfo Set HisProductCode7 = ?,HisProductCode7Source = ? where ProductInfoID= ?", item.YGCGID, item.YGCGID,
		//			prod.ProductInfoID); db.Error != nil {
		//			tx.Rollback()
		//			return db.Error
		//		}
		//		*context += fmt.Sprintf("集采产品ID:HisProductCode7(%s)变更为(%s);", prod.YGCGID, item.YGCGID)
		//	} else if item.YGCGID != prod.HisProductCode7Source {
		//		if db := tx.Exec("Update TB_ProductInfo Set HisProductCode7 = ?, HisProductCode7Source = ? ,HisProductCode7Status = '1' where ProductInfoID= ?", item.YGCGID, item.YGCGID,
		//			prod.ProductInfoID); db.Error != nil {
		//			tx.Rollback()
		//			return db.Error
		//		}
		//		*context += fmt.Sprintf("集采产品ID:HisProductCode7Source(%s)变更为(%s);", prod.HisProductCode7Source, item.YGCGID)
		//	}
		//
	}
	return nil
}

// UpdateTradeCodeInfo 更新流水号
func UpdateTradeCodeInfo(tx *gorm.DB, item *model.ChangeInfoElement, prod model.ProductInfo, context *string) error {
	if item.TradeCode != "" {
		if prod.TradeCode != "" && item.TradeCode != prod.TradeCode {
			if db := tx.Exec("Update TB_TenderCode Set TenderCode =?,UpdateTime = GETDATE() where ProductInfoID = ?",
				item.TradeCode, prod.ProductInfoID); db.Error != nil {
				tx.Rollback()
				return db.Error
			}
			*context += fmt.Sprintf("TradeCode流水号(%s)变更为%s;", prod.TradeCode, item.TradeCode)
		}
		// 插入记录
		if prod.TradeCode == "" {
			if db := tx.Exec("Insert into TB_TenderCode(productinfoid, tendercode, medinsname, medicaretype, startdate, enddate, isvoid, status, tenderstatus) values (?,?,?,?,?,?,?,?,?)",
				prod.ProductInfoID, item.TradeCode, "省标", 1, time.Now(), time.Now(), 0, 0, 0); db.Error != nil {
				tx.Rollback()
				return db.Error
			}
			*context += fmt.Sprintf("产品ID(%v)插入流水号(%s)", prod.ProductInfoID, item.TradeCode)
		}
	}
	return nil
}

// UpdateMedicareCodeInfo2 更新医保代码2.0
func UpdateMedicareCodeInfo2(tx *gorm.DB, item *model.ChangeInfoElement, prod model.ProductInfo, context *string) error {
	if item.MedicareCode != "" {
		// 查询数据库
		var M []model.ProdMedicareCode
		clientDb.DB.Raw("select ChargeRuleID,ProductInfoID,MedicareCode,MedicareCodeTemp ,MedicareCodeStatus from TB_ProductChargeRule where IsVoid = 0 and MedicareType = 1 and ProductInfoID = ?", prod.ProductInfoID).Find(&M)
		if len(M) > 1 {
			tx.Rollback()
			return fmt.Errorf("院内代码:%s,医保代码存在多条,检查院内代码是否对应多条产品ID", item.Code)
		} else if len(M) == 0 { // 无记录
			// 按林老师要求医保代码无记录的直接插入已审核的医保代码
			var insertSql = `INSERT INTO TB_ProductChargeRule ([ProductInfoID],[PriceCharged],[MedicareCode],[MedInsName]
			,[MedicareType],[RepayFlag],[RepayRatio],[AddFeeFlag],[AddRatioFee],[MedicareCodeTemp],[MedicareCodeStatus],[IsVoid])
			select ProductInfoID,ChargePrice,  ? ,N'城镇医保', '1', '0', 0, '0', 0, null, '1','0'
			From TB_ProductInfo where ProductInfoID = ?`
			if db := tx.Exec(insertSql, item.MedicareCode, prod.ProductInfoID); db.Error != nil {
				tx.Rollback()
				return db.Error
			}
			*context += fmt.Sprintf("产品ID(%v)插入医保代码(%s)", prod.ProductInfoID, item.MedicareCode)
			return nil
		}
		// 2024-12-24 林老师反馈医保代码未修改问题，因两边系统已审核的医保代码一致，接口入参数据中没有，但是已审核医保代码与未审核医保代码不一致
		// 2024-12-24 1. 根据审核状态,已审核的判断已审核的医保代码是否与入参一致，不一致直接修改
		// 2. 未审核的医保代码，判断已审核字段是否与入参一致,一致的直接清空未审核医保代码，将状态修改为已审核。不一致的也清空未审核医保代码并修改为已审核状态
		//if M[0].MedicareCode == item.MedicareCode && M[0].MedicareCodeTemp == item.MedicareCode {
		//	return nil
		//}
		// 有记录直接更新
		switch M[0].MedicareCodeStatus {
		// 已审核医保代码
		case "1":
			// 先判断已审核医保代码是否一致
			if M[0].MedicareCode != item.MedicareCode {
				// 不一致 直接修改
				if db := tx.Exec("Update TB_ProductChargeRule Set MedicareCode = ?  where ChargeRuleID = ?", item.MedicareCode, M[0].ChargeRuleID); db.Error != nil {
					tx.Rollback()
					return db.Error
				}
				*context += fmt.Sprintf("MedicareCode:医保代码(%s)变更为(%s);", M[0].MedicareCode, item.MedicareCode)
			}
		// 未审核医保代码
		case "0":
			if M[0].MedicareCode != item.MedicareCode {
				if db := tx.Exec("Update TB_ProductChargeRule Set MedicareCodeTemp = '',MedicareCode = ?,MedicareCodeStatus = 1 where ChargeRuleID = ?", item.MedicareCode, M[0].ChargeRuleID); db.Error != nil {
					tx.Rollback()
					return db.Error
				}
				//*context += fmt.Sprintf("MedicareCodeTemp:未审核医保代码(%s)直接修改为已审核并变更为(%s);", M[0].MedicareCodeTemp, item.MedicareCode)
				*context += fmt.Sprintf("MedicareCode:医保代码(%s)变更为(%s),审核状态(%s)变更为(%s),未审核医保代码(%s)清空;",
					M[0].MedicareCode, item.MedicareCode, M[0].MedicareCodeStatus, "1", M[0].MedicareCodeTemp)
			}
		}
		//for _, v := range M {
		//	// 医保代码不为空
		//	if v.MedicareCode != "" && !strings.HasPrefix(v.MedicareCode, item.MedicareCode) {
		//		if db := tx.Exec("Update TB_ProductChargeRule Set MedicareCode = ? where ChargeRuleID = ?", item.MedicareCode, v.ChargeRuleID); db.Error != nil {
		//			tx.Rollback()
		//			return db.Error
		//		}
		//		*context += fmt.Sprintf("MedicareCode:医保代码(%s)变更为(%s);", v.MedicareCode, item.MedicareCode)
		//	}
		//	// 临时医保代码不为空
		//	if v.MedicareCodeTemp != "" && !strings.HasPrefix(v.MedicareCodeTemp, item.MedicareCode) {
		//		if db := tx.Exec("Update TB_ProductChargeRule Set MedicareCodeTemp = ? where ChargeRuleID = ?", item.MedicareCode, v.ChargeRuleID); db.Error != nil {
		//			tx.Rollback()
		//			return db.Error
		//		}
		//		*context += fmt.Sprintf("MedicareCodeTemp:未审核医保代码(%s)变更为(%s);", v.MedicareCodeTemp, item.MedicareCode)
		//	}
		//	// 医保代码为空
		//	if v.MedicareCode == "" {
		//		if db := tx.Exec("Update TB_ProductChargeRule set MedicareCode = ? where ChargeRuleID = ?", item.MedicareCode, v.ChargeRuleID); db.Error != nil {
		//			tx.Rollback()
		//			return db.Error
		//		}
		//		*context += fmt.Sprintf("MedicareCode:医保代码(%s)变更为(%s);", v.MedicareCode, item.MedicareCode)
		//	}
		//}
	}
	return nil
}

// 因为医保代码表关系复杂,此方法是根据查询所有产品信息来修改,医保代码可能出现多条等，故弃用此方法,改为2.0,需要时在单独查询数据库
// UpdateMedicareCodeInfo 更新医保代码1.0
//func UpdateMedicareCodeInfo(tx *gorm.DB, item *ChangeInfoElement, prod model.ProductInfo, context *string) error {
//	if *item.MedicareCode != "" {
//		// 如果医保代码和国家医保代码都为空则插入医保代码
//		if prod.MedicareCode == "" && prod.CountryMedicareCode == "" {
//			db := tx.Exec("Insert into TB_ProductChargeRule(productinfoid, pricecharged, medicarecode, medinsname, medicaretype, repayflag, repayratio, addfeeflag, addratiofee, isvoid, medicarecodestatus) values (?,?,?,?,?,?,?,?,?,?,?)",
//				prod.ProductInfoID, prod.ChargePrice, *item.MedicareCode, "城镇医保", 1, 0, 0, 0, 0, 0, 1)
//			if db.Error != nil {
//				tx.Rollback()
//				return db.Error
//			}
//			*context += fmt.Sprintf("产品ID(%v)插入医保代码(%s)", prod.ProductInfoID, *item.MedicareCode)
//		} else {
//			// prod.MedicareCode不为空,且不以item.MedicareCode开头的则修改
//			if prod.MedicareCode != "" && !strings.HasPrefix(prod.MedicareCode, *item.MedicareCode) {
//				db := tx.Exec("Update TB_ProductChargeRule Set MedicareCode = ? where ProductInfoID = ? and MedicareType = 1 ",
//					item.MedicareCode, prod.ProductInfoID)
//				if db.Error != nil {
//					tx.Rollback()
//					return db.Error
//				}
//				*context += fmt.Sprintf("MedicareCode:医保代码(%s)变更为(%s);", prod.MedicareCode, *item.MedicareCode)
//			}
//			if prod.CountryMedicareCode != "" && !strings.HasPrefix(prod.CountryMedicareCode, *item.MedicareCode) {
//				db := tx.Exec("Update TB_ProductChargeRule Set MedicareCode = ? where ProductInfoID = ? ")
//			}
//		}
//	}
//	return nil
//}

// UpdateJCSysInfo 更新集采系统信息
func UpdateJCSysInfo(tx *gorm.DB, item *model.ChangeInfoElement, prod model.ProductInfo, context *string) error {
	if prod.SysCode == "" && prod.SysId == "" { // 如果同时为空则代表无记录
		if item.SysID != "" || item.SysCode != "" {
			if db := tx.Exec("Insert Into TB_ProductInfoJCSysCode(Prod_Id, SysId, SysCode, IsVoid, CreateTime) values (?,?,?,?,getdate())",
				prod.ProductInfoID, item.SysID, item.SysCode, 0); db.Error != nil {
				tx.Rollback()
				return db.Error
			}
			*context += fmt.Sprintf("插入集采系统信息:Sysid(%s),SysCode(%s)", item.SysID, item.SysCode)
		}
	} else {
		if item.SysID != "" && item.SysID != prod.SysId {
			if db := tx.Exec("Update TB_ProductInfoJCSysCode set SysId = ? where Prod_Id =?", item.SysID,
				prod.ProductInfoID); db.Error != nil {
				tx.Rollback()
				return db.Error
			}
			*context += fmt.Sprintf("更新集采系统信息:SysId(%s)变更为(%s);", prod.SysId, item.SysID)
		}
		if item.SysCode != "" && item.SysCode != prod.SysCode {
			if db := tx.Exec("Update TB_ProductInfoJCSysCode set SysCode = ? where Prod_Id =?", item.SysCode,
				prod.ProductInfoID); db.Error != nil {
				tx.Rollback()
				return db.Error
			}
			*context += fmt.Sprintf("更新集采系统信息:SysCode(%s)变更为(%s)", prod.SysCode, item.SysCode)
		}
	}
	return nil
}

// UpdateProductSupplyStatus 停供
func UpdateProductSupplyStatus(tx *gorm.DB, prod model.ProductInfo, context *string) error {
	// 1. 查询白名单表
	var row int
	var sql = `select count(1) from TB_DepartmentApply where ProductInfoID = ? `
	clientDb.DB.Raw(sql, prod.ProductInfoID).Scan(&row)
	if row == 0 {
		return nil
	}
	// 修改白名单
	var sql2 = `update TB_DepartmentApply set IsVoid = 1,UpdateTime = getdate() where ProductInfoID = ? and IsVoid = 0`
	if db := tx.Exec(sql2, prod.ProductInfoID); db.Error != nil {
		tx.Rollback()
		return db.Error
	}
	*context += fmt.Sprintf("停供产品处理更新产品白名单表产品ID:%v,", prod.ProductInfoID)
	return nil
}
