package controller

import (
	clientDb "WorkloadQuery/db"
	"WorkloadQuery/model"
)

type QueryAccountEntryTime struct {
	StartTime string `json:"startTime" binding:"required"`
	EndTime   string `json:"endTime" binding:"required"`
}

func (ae *QueryAccountEntryTime) NoAccountEntryQuery() (noAccountEntryBills *[]model.NoAccountEntry) {
	clientDb.DB.Raw(clientDb.NoAccountEntrySql, ae.StartTime, ae.EndTime).Find(&noAccountEntryBills)
	for i := range *noAccountEntryBills {
		// tempTime, _ := time.Parse("2006-01-02T15:04:05Z", (*noAccountEntryBills)[i].CreateTime)
		// (*noAccountEntryBills)[i].CreateTime = tempTime.Format("2006-01-02 15:04:05")
		// 序号
		(*noAccountEntryBills)[i].Index = i + 1
	}
	return
}
