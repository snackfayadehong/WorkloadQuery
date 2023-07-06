package model

// NoAccountEntry 调拨未上账单据模型
type NoAccountEntry struct {
	Index                  int    // 序号
	DepartmentCollarCode   string `gorm:"column:DepartmentCollarCode"` // 单号
	BLDate                 string `gorm:"BLDate"`                      // 库房审核时间
	LeadingDepartmentName  string `gorm:"LeadingDepartmentName"`       // 领用科室
	LeaderName             string `gorm:"LeaderName"`                  // 领用人
	TreasuryDepartmentName string `gorm:"TreasuryDepartmentName"`      // 出库科室
	BLMakerName            string `gorm:"BLMakerName"`                 // 制单人
	Flag                   string `gorm:"Flag"`                        // 状态
}
