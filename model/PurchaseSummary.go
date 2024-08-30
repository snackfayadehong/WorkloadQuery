package model

// 库房廖小凤反馈统计为配送采购订单及未到货采购订单明细并可导出 --2023-09-12

// PurchaseSummary 采购订单主表Model
type PurchaseSummary struct {
	PurchaseSummaryID   string  `gorm:"column:PurchaseSummaryID"`   // 订单ID
	PurchaseSummaryCode string  `gorm:"column:PurchaseSummaryCode"` // 订单号
	AllMoney            float32 `gorm:"column:AllMoney"`            // 订单金额
	DepartmentName      string  `gorm:"column:DepartmentName"`      // 入库科室
	StatusName          string  `gorm:"column:StatusName"`          // 订单状态
	MakeName            string  `gorm:"column:MakeName"`            // 开单人员
	AuditorDate         string  `gorm:"column:AuditorDate"`         // 订单审核时间
	GoodStatusName      string  `gorm:"column:GoodStatusName"`      // 订单物流状态
	SupplierName        string  `gorm:"column:SupplierName"`        // 供货公司
	Remark              string  `gorm:"column:Remark"`              // 订单备注
}

// PurchaseSummaryDetail 采购订单明细Model
type PurchaseSummaryDetail struct {
	PurchaseSummaryID string  `gorm:"column:PurchaseSummaryID"`  // 订单ID
	Code              string  `gorm:"column:Code"`               // 院内代码
	ProdName          string  `gorm:"column:ProdName"`           // 产品名称
	HospitalSpec      string  `gorm:"column:HospitalSpec"`       // 院内规格
	EnterpriseName    string  `gorm:"column:EnterpriseName"`     // 供货公司
	UnitPrice         float32 `gorm:"column:UnitPrice"`          // 单价
	SpecName          string  `gorm:"column:SpecName"`           // 规格
	Amount            float32 `gorm:"column:Amount"`             // 数量
	FactInAmount      float32 `gorm:"column:FactInAmount"`       // 到货数量
	RefundAmount      float32 `gorm:"column:RefundAmount"`       // 关闭数量
	NotDelivered      float32 `gorm:"column:NotDeliveredAmount"` // 未到货数量
	Remark            string  `gorm:"column:Remark"`             // 明细备注
}

type AllPurchaseSummary struct {
	PurchaseSummary
	Children *[]PurchaseSummaryDetail `gorm:"_" json:"children"`
}
