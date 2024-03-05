package model

// ProductInfo 查询数据库产品字典信息
type ProductInfo struct {
	ProductInfoID         int    `gorm:"column:ProductInfoID"`         // 产品ID
	Code                  string `gorm:"column:Code"`                  // 院内代码
	HospitalName          string `gorm:"column:HospitalName"`          // 院内名称
	HospitalSpec          string `gorm:"column:HospitalSpec"`          // 院内规格
	YGCGID                string `gorm:"column:YGCGID"`                // 网采平台产品ID
	HisProductCode7Source string `gorm:"column:HisProductCode7Source"` // 集采平台产品临时ID 审核过后会改变HisProductCode7值
	HisProductCode7Status int    `gorm:"column:HisProductCode7Status"` // 集采平台产品ID状态
	CusCategoryCode       string `gorm:"column:CusCategoryCode"`       // 104分类编码
	ParentCusCategoryCode string `gorm:"column:ParentCusCategoryCode"` // 104分类编码(第3级)
	TradeCode             string `gorm:"column:TradeCode"`             // 交易编码
	ChargePrice           string `gorm:"column:ChargePrice"`           // 收费价格
	MedicareCode          string `gorm:"column:MedicareCode"`          // 医保编码
	SysCode               string `gorm:"column:SysCode"`               // 集采系统编码
	SysId                 string `gorm:"column:SysId"`                 // 集采系统编号
	IsVoid                int    `gorm:"column:IsVoid"`                // 0：启用 1：停用
	PurState              int    `gorm:"column:PurState"`              // 0：供货 1：停止供货
}

type GetProduct interface {
	GetProductInfo(Where []string) (*[]ProductInfo, error)
}
type ChangeProduct interface {
	ChangeProductInfo(*[]ProductInfo) error
}
