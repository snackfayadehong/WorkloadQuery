package controller

import (
	clientDb "WorkloadQuery/db"
	"WorkloadQuery/model"
	"time"
)

func NoAccountEntryQuery(endTime *string) (bills *[]model.NoAccountEntry) {
	clientDb.DB.Raw(clientDb.NoAccountEntrySql, endTime).Find(&bills)
	for i, _ := range *bills {
		tempTime, _ := time.Parse("2006-01-02T15:04:05Z", (*bills)[i].BLDate)
		(*bills)[i].BLDate = tempTime.Format("2006-01-02 15:04:05")
	}
	return
}
