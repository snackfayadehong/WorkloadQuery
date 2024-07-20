package controller

import (
	clientDb "WorkloadQuery/db"
	"WorkloadQuery/model"
)

type QueryWorkloadTime struct {
	StartTime string `json:"startTime" binding:"required"`
	EndTime   string `json:"endTime" binding:"required"`
}

// // UserWorkloadQuery  工作量查询（2024-07-17 加入普耗数据后停用,有UserWorkloadQuery_New代替）
//
//	func (workload *QueryWorkloadTime) UserWorkloadQuery() *[]model.UserWorkloadInfo {
//		var ProdAccept []model.UserProdAccept                      // 入库
//		var DpProd []model.DepartmentCollar                        // 出库
//		var RefProd []model.RefundProd                             // 退货
//		UserWorkloadMap := make(map[string]model.UserWorkloadInfo) // 合并数据map
//		var UserWorkload []model.UserWorkloadInfo                  // 合并后的数据切片
//		// 入库
//		clientDb.DB.Raw(clientDb.UserProdAcceptSql, workload.StartTime, workload.EndTime, workload.StartTime, workload.EndTime).Find(&ProdAccept)
//		// 出库
//		clientDb.DB.Raw(clientDb.UserProdDpcSql, workload.StartTime, workload.EndTime, workload.StartTime, workload.EndTime).Find(&DpProd)
//		// 退货
//		clientDb.DB.Raw(clientDb.UserProdRefundSql, workload.StartTime, workload.EndTime, workload.StartTime, workload.EndTime).Find(&RefProd)
//
//		// 合并数据
//		for i := 0; i < len(ProdAccept) || i < len(DpProd) || i < len(RefProd); i++ {
//			if i < len(ProdAccept) {
//				usm := UserWorkloadMap[ProdAccept[i].Name] // 找到对应Name的map并把value赋值给usm
//				usm.Name = ProdAccept[i].Name
//				usm.ProdAcBill = ProdAccept[i].ProdAcBill
//				usm.ProdAcSpec = ProdAccept[i].ProdAcSpec
//				usm.ProdAcTotal = ProdAccept[i].ProdAcTotal
//				UserWorkloadMap[ProdAccept[i].Name] = usm
//			}
//			if i < len(DpProd) {
//				usm := UserWorkloadMap[DpProd[i].Name]
//				usm.Name = DpProd[i].Name
//				usm.ProdDpBill = DpProd[i].ProdDpBill
//				usm.ProdDpSpec = DpProd[i].ProdDpSpec
//				usm.ProdDpTotal = DpProd[i].ProdDpTotal
//				UserWorkloadMap[DpProd[i].Name] = usm
//			}
//			if i < len(RefProd) {
//				usm := UserWorkloadMap[RefProd[i].Name]
//				usm.Name = RefProd[i].Name
//				usm.RefBill = RefProd[i].RefBill
//				usm.RefSpec = RefProd[i].RefSpec
//				usm.RefTotal = RefProd[i].RefTotal
//				UserWorkloadMap[RefProd[i].Name] = usm
//			}
//		}
//		// 将合并后的数据写入切片
//		for _, v := range UserWorkloadMap {
//			UserWorkload = append(UserWorkload, v)
//		}
//		return &UserWorkload
//	}

// UserWorkloadQuery_New 2014-07-17 工作量查询逻辑执行存储过程
func (workload *QueryWorkloadTime) UserWorkloadQuery_New() *[]model.UserWorkloadInfoNew {
	// 查询工作量
	var w []model.UserWorkloadInfoNew
	//var sql = `EXEC dbo.WorkloadReport ?, ?,'',''`
	//clientDb.DB.Raw(sql, workload.StartTime, workload.EndTime).Find(&w)
	var sql = `EXEC dbo.WorkloadReport ?,?`
	db := clientDb.DB.Raw(sql, workload.StartTime, workload.EndTime).Find(&w)
	if db.Error != nil {
		return nil
	}
	//rows, _ := clientDb.DB.Raw(sql, workload.StartTime, workload.EndTime).Rows()
	//defer rows.Close()
	//for rows.Next() {
	//	clientDb.DB.ScanRows(rows, &w)
	//}
	// 计算每个人的总数据
	for i, v := range w {
		w[i].TotalBillAmount = v.In_BillNum + v.Out_BillNum + v.Back_BillNum
		w[i].TotalSpecAmount = v.In_ProdSpecNum + v.Out_ProdSpecNum + v.Back_ProdSpecNum
		w[i].TotalTotalAmount = v.In_TotalAmount + v.Out_TotalAmount + v.Back_TotalAmount
	}
	return &w
}
