package controller

import (
	clientDb "WorkloadQuery/db"
	"WorkloadQuery/model"
	"time"
)

type QueryUncheckTime struct {
	StartTime string `json:"startTime" binding:"required"`
	EndTime   string `json:"endTime" binding:"required"`
}

// UnCheckBillQuery 计费未核对数据查询
func (ub *QueryUncheckTime) UnCheckBillQuery() (uncheckBills *[]model.UnCheckBill) {
	clientDb.DB.Raw(clientDb.UnCheckBillSql, ub.StartTime, ub.EndTime).Find(&uncheckBills)
	for i := range *uncheckBills {
		tempTime, _ := time.Parse("2006-01-02T15:04:05Z", (*uncheckBills)[i].CheckDate)
		(*uncheckBills)[i].CheckDate = tempTime.Format("2006-01-02 15:04:05")
		// 序号
		(*uncheckBills)[i].Index = i + 1
	}
	return
}
