package controller

import (
	clientDb "WorkloadQuery/db"
	"WorkloadQuery/model"
	"fmt"
	"gorm.io/gorm"
	"time"
)

/*
	根据老物资采购系统传过来的字典信息变更
*/

// ChangeInfoElement 字典信息入参
type ChangeInfoElement struct {
	Code         string  `json:"Code"`                   // 院内代码
	HospitalName *string `json:"HospitalName,omitempty"` // 院内产品名称
	HospitalSpec *string `json:"HospitalSpec,omitempty"` // 院内规格
	YGCGID       *string `json:"YGCGID,omitempty"`       // 网采平台产品ID
	TradeCode    *string `json:"TradeCode,omitempty"`    // 商品代码，商品代码或挂网流水号（平台供货商填写交易编码的内容）
	MedicareCode *string `json:"MedicareCode,omitempty"` // 医保代码
	CategoryCode *string `json:"CategoryCode,omitempty"` // 18类分类代码，匹配104分类三级目录,如果匹配成功不修改，否则修改为三级目录
	SysCode      *string `json:"SysCode,omitempty"`      // 系统编码
	SysID        *string `json:"SysID,omitempty"`        // 系统编号
}
type RequestInfo struct {
	C []ChangeInfoElement
}

const UpdateCateCodeSql = "Update TB_ProductInfo Set CusCategoryCode = ? where ProductInfoID = ?"
const UpdateHospitalSpecSql = "Update TB_ProductInfo set HospitalSpec = ? where ProductInfoID = ?"
const UpdateHospitalNameSql = "Update TB_ProductInfo Set HisProductCode3 =? where ProductInfoID = ?"

/*
ChangeProductInfo
更改产品基本信息
*/
func (i *RequestInfo) ChangeProductInfo(prod *[]model.ProductInfo, ip string) error {
	// 开启事务
	tx := clientDb.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// 查找入参和物资系统返回查询结果中相同的记录 时间复杂度O(M+N)
	// 使用Map存储prod切片中的Code
	pMap := make(map[string]int, len(*prod))
	for index, v := range i.C {
		pMap[v.Code] = index
	}
	// 在入参中找到Code相同的
	for _, item := range i.C {
		if pIndex, ok := pMap[item.Code]; ok {
			var updateMsg string
			// 修改字典信息的业务逻辑
			// 1. 104分类；
			err := UpdateCategoryCode(tx, &item, (*prod)[pIndex], &updateMsg)
			if err != nil {
				return err
			}
			// 2. 更新院内产品名称、规格信息
			err = UpdateHospitalInfo(tx, &item, (*prod)[pIndex], &updateMsg)
			if err != nil {
				return err
			}
			// 3. 判断集采审核状态并更新集采信息
			err = UpdateYgcgidInfo(tx, &item, (*prod)[pIndex], &updateMsg)
			if err != nil {
				return err
			}
			// 4. 更新TradeCode流水号
			err = UpdateTradeCodeInfo(tx, &item, (*prod)[pIndex], &updateMsg)
			if err != nil {
				return err
			}
			// 5. 更新医保代码
			err = UpdateMedicareCodeInfo(tx, &item, (*prod)[pIndex], &updateMsg)
			if err != nil {
				return err
			}
			// 6. 写入系统编码系统编号
			err = UpdateJCSysInfo(tx, &item, (*prod)[pIndex], &updateMsg)
			if err != nil {
				return err
			}
			// 写入日志表
			if updateMsg != "" {
				db := tx.Exec("Insert Into TB_PoductInfoChangeLog(Prod_Id,Context,IP,UpdateTime) values(?,?,?,?)",
					(*prod)[pIndex].ProductInfoID, updateMsg, ip, time.Now())
				if db.Error != nil {
					tx.Rollback()
					return db.Error
				}
			}
		}
	}
	tx.Commit()
	return nil
}

/*
GetProductInfo
获取物资产品字典信息,返回不重复的字典信息
string返回记录哪些字典信息是重复的
*/
func (i *RequestInfo) GetProductInfo(Where []string) (*[]model.ProductInfo, error) {
	var prod *[]model.ProductInfo         // 原始记录
	var NoRepeatProd []model.ProductInfo  // 保留不重复的记录
	var repeatMap = make(map[string]bool) // 重复记录
	var msg string
	db := clientDb.DB.Raw(clientDb.QueryProd, Where).Find(&prod)
	if db.Error != nil {
		return prod, db.Error
	}
	// 检查 查询结果中同一院内编码是否存在多条记录,且非停用或停供产品
	seen := make(map[string]bool)
	for _, el := range *prod {
		// 非停用或者停供且不重复的记录添加到map,和切片中
		if !seen[el.Code] && el.PurState == 0 && el.IsVoid == 0 {
			seen[el.Code] = true
			NoRepeatProd = append(NoRepeatProd, el)
			continue
		}
		// 记录哪些是重复的,只记录一次
		if !repeatMap[el.Code] {
			repeatMap[el.Code] = true
		}
	}
	for key := range repeatMap {
		msg += fmt.Sprintf("%s有重复字典记录或供货关系异常;", key)
	}
	if msg == "" {
		return &NoRepeatProd, nil
	}
	return &NoRepeatProd, fmt.Errorf(msg)
}

// UpdateCategoryCode 更新104分类
func UpdateCategoryCode(tx *gorm.DB, item *ChangeInfoElement, prod model.ProductInfo, context *string) error {
	// 入参中104分类为第3级,物资系统为第4级,查询物资系统第4级代码对应的第3级与入参是否相同
	// 如果不相同则修改物资为第3级,相同则不更新
	if *item.CategoryCode != prod.ParentCusCategoryCode {
		// 修改为第3级
		db := tx.Exec(UpdateCateCodeSql, item.CategoryCode, prod.ProductInfoID)
		if db.Error != nil {
			tx.Rollback()
			return db.Error
		}
		*context = fmt.Sprintf("产品id:%v,CusCategoryCode(%s)变更为(%s);", prod.ProductInfoID, prod.CusCategoryCode, *item.CategoryCode)
	}
	return nil
}

// UpdateHospitalInfo 更新院内名称、院内规格
func UpdateHospitalInfo(tx *gorm.DB, item *ChangeInfoElement, prod model.ProductInfo, context *string) error {
	if *item.HospitalSpec != "" && *item.HospitalSpec != prod.HospitalSpec {
		db := tx.Exec(UpdateHospitalSpecSql, item.HospitalSpec, prod.ProductInfoID)
		if db.Error != nil {
			tx.Rollback()
			return db.Error
		}

	}
	if *item.HospitalName != "" && *item.HospitalName != prod.HospitalName {
		db := tx.Exec(UpdateHospitalNameSql, item.HospitalName, prod.ProductInfoID)
		if db.Error != nil {
			tx.Rollback()
			return db.Error
		}
		*context += fmt.Sprintf("HospitalSpec:院内规格(%s)变更为(%s);HospitalName:院内名称(%s)变更为(%s);", prod.HospitalSpec, *item.HospitalSpec,
			prod.HospitalName, *item.HospitalName)
	}
	return nil
}

// UpdateYgcgidInfo 更新网采平台产品ID
func UpdateYgcgidInfo(tx *gorm.DB, item *ChangeInfoElement, prod model.ProductInfo, context *string) error {
	// 1 已审核  null '' 0 为未审核
	if *item.YGCGID != "" {
		if prod.HisProductCode7Status == 1 && *item.YGCGID != prod.YGCGID {
			db := tx.Exec("Update TB_ProductInfo Set HisProductCode7 = ? where ProductInfoID= ?", item.YGCGID,
				prod.ProductInfoID)
			if db.Error != nil {
				tx.Rollback()
				return db.Error
			}
			*context += fmt.Sprintf("集采产品ID:HisProductCode7(%s)变更为(%s);", prod.YGCGID, *item.YGCGID)
		} else if *item.YGCGID != prod.HisProductCode7Source {
			db := tx.Exec("Update TB_ProductInfo Set HisProductCode7Source = ? where ProductInfoID= ?", item.YGCGID,
				prod.ProductInfoID)
			if db.Error != nil {
				tx.Rollback()
				return db.Error
			}
			*context += fmt.Sprintf("集采产品ID:HisProductCode7Source(%s)变更为(%s);", prod.HisProductCode7Source, *item.YGCGID)
		}

	}
	return nil
}

// UpdateTradeCodeInfo 更新流水号
func UpdateTradeCodeInfo(tx *gorm.DB, item *ChangeInfoElement, prod model.ProductInfo, context *string) error {
	if *item.TradeCode != "" && *item.TradeCode != prod.TradeCode {
		db := tx.Exec("Update TB_TenderCode Set TenderCode =?,UpdateTime = GETDATE() where ProductInfoID = ?",
			item.TradeCode, prod.ProductInfoID)
		if db.Error != nil {
			tx.Rollback()
			return db.Error
		}
		*context += fmt.Sprintf("TradeCode流水号(%s)变更为%s;", prod.TradeCode, *item.TradeCode)
	}
	return nil
}

// UpdateMedicareCodeInfo 更新医保代码
func UpdateMedicareCodeInfo(tx *gorm.DB, item *ChangeInfoElement, prod model.ProductInfo, context *string) error {
	if *item.MedicareCode != "" && *item.MedicareCode != prod.MedicareCode {
		db := tx.Exec("Update TB_ProductChargeRule Set MedicareCode = ? where ProductInfoID = ?",
			item.MedicareCode, prod.ProductInfoID)
		if db.Error != nil {
			tx.Rollback()
			return db.Error
		}
		*context += fmt.Sprintf("MedicareCode:医保代码(%s)变更为(%s);", prod.MedicareCode, *item.MedicareCode)
	}
	return nil
}

// UpdateJCSysInfo 更新集采系统信息
func UpdateJCSysInfo(tx *gorm.DB, item *ChangeInfoElement, prod model.ProductInfo, context *string) error {
	if prod.SysCode == "" && prod.SysId == "" { // 如果同时为空则代表无记录
		if *item.SysID != "" || *item.SysCode != "" {
			db := tx.Exec("Insert Into TB_ProductInfoJCSysCode(Prod_Id, SysId, SysCode, IsVoid, CreateTime) values (?,?,?,?,getdate())",
				prod.ProductInfoID, item.SysID, item.SysCode, 0)
			if db.Error != nil {
				tx.Rollback()
				return db.Error
			}
			*context += fmt.Sprintf("插入集采系统信息:Sysid(%s),SysCode(%s)", *item.SysID, *item.SysCode)
		}
	} else {
		if *item.SysID != "" && *item.SysID != prod.SysId {
			db := tx.Exec("Update TB_ProductInfoJCSysCode set SysId = ? where Prod_Id =?", item.SysID,
				prod.ProductInfoID)
			if db.Error != nil {
				tx.Rollback()
				return db.Error
			}
			*context += fmt.Sprintf("更新集采系统信息:SysId(%s)变更为(%s);", prod.SysId, *item.SysID)
		}
		if *item.SysCode != "" && *item.SysCode != prod.SysCode {
			db := tx.Exec("Update TB_ProductInfoJCSysCode set SysCode = ? where Prod_Id =?", item.SysCode,
				prod.ProductInfoID)
			if db.Error != nil {
				tx.Rollback()
				return db.Error
			}
			*context += fmt.Sprintf("更新集采系统信息:SysCode(%s)变更为(%s)", prod.SysCode, *item.SysCode)
		}
	}
	return nil
}
