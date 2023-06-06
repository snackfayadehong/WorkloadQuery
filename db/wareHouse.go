package clientDb

// UserProdAccept 验收
type UserProdAccept struct {
	Name        string `gorm:"column:MEnName"`
	ProdAcSpec  int    `gorm:"column:prodSpecNum"`
	ProdAcBill  int    `gorm:"column:billNum"`
	ProdAcTotal string `gorm:"column:totalAmount"`
}

// DepartmentCollar 出库
type DepartmentCollar struct {
	Name        string `gorm:"column:BLMakerName"`
	ProdDpSpec  int    `gorm:"column:DpSpecNum"`
	ProdDpBill  int    `gorm:"column:DpBillNum"`
	ProdDpTotal string `gorm:"column:DpTotalAmount"`
}

// RefundProd 退货
type RefundProd struct {
	Name     string `gorm:"column:EmployeeName"`
	RefSpec  int    `gorm:"column:ReFSpecNum"`
	RefBill  int    `gorm:"column:RefBillNum"`
	RefTotal string `gorm:"column:RefTotalAmount"`
}
type UserWorkloadInfo struct {
	Name             string `json:"name"`
	ProdAcSpec       int    `json:"prodAcSpec"`
	ProdAcBill       int    `json:"prodAcBill"`
	ProdAcTotal      string `json:"prodAcTotal"`
	ProdDpSpec       int    `json:"prodDpSpec"`
	ProdDpBill       int    `json:"prodDpBill"`
	ProdDpTotal      string `json:"prodDpTotal"`
	RefSpec          int    `json:"refSpec"`
	RefBill          int    `json:"refBill"`
	RefTotal         string `json:"refTotal"`
	TotalBillAmount  int    `json:"total_bill_amount"`
	TotalSpecAmount  int    `json:"total_spec_amount"`
	TotalTotalAmount string `json:"total_total_amount"`
}

// UserWorkloadQuery  工作量查询
func UserWorkloadQuery(startTime string, endTime string) []UserWorkloadInfo {
	var ProdAccept []UserProdAccept                      // 入库
	var DpProd []DepartmentCollar                        // 出库
	var RefProd []RefundProd                             // 退货
	UserWorkloadMap := make(map[string]UserWorkloadInfo) // 合并数据map
	var UserWorkload []UserWorkloadInfo                  // 合并后的数据切片
	// 入库
	DB.Raw(UserProdAcceptSql, startTime, endTime).Find(&ProdAccept)
	// 出库
	DB.Raw(UserProdDpcSql, startTime, endTime).Find(&DpProd)
	// 退货
	DB.Raw(UserProdRefundSql, startTime, endTime).Find(&RefProd)

	// 合并数据
	for i := 0; i < len(ProdAccept) || i < len(DpProd) || i < len(RefProd); i++ {
		if i < len(ProdAccept) {
			usm := UserWorkloadMap[ProdAccept[i].Name] // 找到对应Name的map并把value赋值给usm
			usm.Name = ProdAccept[i].Name
			usm.ProdAcBill = ProdAccept[i].ProdAcBill
			usm.ProdAcSpec = ProdAccept[i].ProdAcSpec
			usm.ProdAcTotal = ProdAccept[i].ProdAcTotal
			UserWorkloadMap[ProdAccept[i].Name] = usm
		}
		if i < len(DpProd) {
			usm := UserWorkloadMap[DpProd[i].Name]
			usm.Name = DpProd[i].Name
			usm.ProdDpBill = DpProd[i].ProdDpBill
			usm.ProdDpSpec = DpProd[i].ProdDpSpec
			usm.ProdDpTotal = DpProd[i].ProdDpTotal
			UserWorkloadMap[DpProd[i].Name] = usm
		}
		if i < len(RefProd) {
			usm := UserWorkloadMap[RefProd[i].Name]
			usm.Name = RefProd[i].Name
			usm.RefBill = RefProd[i].RefBill
			usm.RefSpec = RefProd[i].RefBill
			usm.RefTotal = RefProd[i].RefTotal
			UserWorkloadMap[RefProd[i].Name] = usm
		}
	}
	// 将合并后的数据写入切片
	for _, v := range UserWorkloadMap {
		UserWorkload = append(UserWorkload, v)
	}
	return UserWorkload
}
