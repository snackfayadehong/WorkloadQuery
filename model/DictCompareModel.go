package model

// LocalDictRow 本地数据库字典数据结构
type LocalDictRow struct {
	Ypdm          string  `gorm:"column:ypdm" json:"ypdm"` // 别名 a.Code
	ProductInfoID string  `gorm:"column:ProductInfoID" json:"ProductInfoID"`
	Ypmc          string  `gorm:"column:ypmc" json:"ypmc"`   // 通用名
	Ypgg          string  `gorm:"column:ypgg" json:"ypgg"`   // 规格
	Kfdw          string  `gorm:"column:kfdw" json:"kfdw"`   // 单位
	Kfcgj         float64 `gorm:"column:kfcgj" json:"kfcgj"` // 采购价
	Kflsj         float64 `gorm:"column:kflsj" json:"kflsj"` // 零售价
	Kfdm          string  `gorm:"column:kfdm" json:"kfdm"`   // 库房代码
	Ghdw          string  `gorm:"column:ghdw" json:"ghdw"`   // 供货单位
	Gsdm          string  `gorm:"column:gsdm" json:"gsdm"`   // 公司代码
}

// DictCompareResult 前端展示的单项对比结果
type DictCompareResult struct {
	Label      string      `json:"label"`      // 中文名称
	Field      string      `json:"field"`      // 字段名
	LocalValue interface{} `json:"localValue"` // 本地值
	HisValue   interface{} `json:"hisValue"`   // HIS值
	IsMatch    bool        `json:"isMatch"`    // 是否一致
}

// HisResponseData HIS接口返回的最外层结构
type HisResponseData struct {
	Users []map[string]interface{} `json:"users"`
}
