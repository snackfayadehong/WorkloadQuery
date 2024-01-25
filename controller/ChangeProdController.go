package controller

import (
	clientDb "WorkloadQuery/db"
	"WorkloadQuery/model"
	"fmt"
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

func (i *ChangeInfoElement) ChangeProductInfo() {

}

/*
GetProductInfo
返回不重复的字典信息
string返回记录哪些字典信息是重复的
*/
func (i *ChangeInfoElement) GetProductInfo(Code string) (*[]model.ProductInfo, string) {
	var prod *[]model.ProductInfo        // 原始记录
	var NoRepeatProd []model.ProductInfo // 保留不重复的记录
	var msg string                       // 返回重复记录信息
	clientDb.DB.Raw(clientDb.QueryProd, Code).Find(&prod)
	// 检查 查询结果中同一院内编码是否存在多条记录
	seen := make(map[string]bool)
	for _, el := range *prod {
		// 不重复将当前记录添加到Map,并添加到新切片
		if !seen[el.Code] {
			seen[el.Code] = true
			NoRepeatProd = append(NoRepeatProd, el)
		}
		msg += fmt.Sprintf("%s存在多条字典信息", el.Code)
	}
	return &NoRepeatProd, msg
}
