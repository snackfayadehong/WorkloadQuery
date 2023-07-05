package model

// UnCheckBill 计费未核对数据模型
type UnCheckBill struct {
	Index       int    // 序号
	BillNo      string `gorm:"column:BillNo"`      // 单据号
	PatName     string `gorm:"column:PatName"`     // 病人名称
	Section     string `gorm:"column:Section"`     // 执行科室
	Doctor      string `gorm:"column:Doctor"`      // 医生
	OperateName string `gorm:"column:OperateName"` // 操作人
	CheckDate   string `gorm:"column:CheckDate"`   // 审核时间
	CheckStatus string `gorm:"column:Status"`      // 状态
}
