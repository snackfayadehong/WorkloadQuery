package model

// NoAccountEntry 调拨未上账单据模型
type NoAccountEntry struct {
	Index                 int    // 序号
	DepartmentCollarCode  string `gorm:"column:DepartmentCollarCode"`  // 单号
	LeadingDepartmentName string `gorm:"column:LeadingDepartmentName"` // 领用科室
	LeaderName            string `gorm:"LeaderName"`                   // 领用人
	BLMakerName           string `gorm:"BLMakerName"`                  // 制单人
	StatisticalTime       string `gorm:"column:StatisticalTime"`       // 数据统计时间
	CreateTime            string `gorm:"column:CreateTime"`            // 总库出库时间
	Code                  string `gorm:"Code"`                         // 材料代码
	ProdName              string `gorm:"column:ProdName"`              // 材料名称
	SpecModelName         string `gorm:"column:SpecModelName"`         // 型号|规格
	UnitName              string `gorm:"column:UnitName"`              // 单位
	ChargePrice           string `gorm:"column:ChargePrice"`           // 单价
	Amount                string `gorm:"column:amount"`                // 数量
}
