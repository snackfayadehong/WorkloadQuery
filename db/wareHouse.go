package clientDb

import (
	"encoding/json"
)

type User struct {
	UserName string `gorm:"column:EmployeeName"`
}

// UserProdAccept 验收
type UserProdAccept struct {
	Name        string `gorm:"column:Oper"`
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

// UserWorkloadQuery  工作量查询
func UserWorkloadQuery(startTime string, endTime string) (Accept string, Outbound string, Refund string) {
	var ProdAccept []UserProdAccept // 入库
	var DpProd []DepartmentCollar   // 出库
	var RefProd []RefundProd        // 退货
	// 入库
	DB.Raw(UserProdAcceptSql, startTime, endTime).Find(&ProdAccept)
	// 出库
	DB.Raw(UserProdDpcSql, startTime, endTime).Find(&DpProd)
	// 退货
	DB.Raw(UserProdRefundSql, startTime, endTime).Find(&RefProd)
	_, err := json.Marshal(ProdAccept)
	if err != nil {
		return
	}
}
