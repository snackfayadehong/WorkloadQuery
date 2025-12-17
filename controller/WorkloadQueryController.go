package controller

import (
	clientDb "WorkloadQuery/db"
	"WorkloadQuery/model"
	"fmt"
)

type QueryDate model.WorkloadRequest

func (q QueryDate) GetWorkloadStatistics() (*[]model.WorkloadData, error) {
	var d []model.WorkloadData
	var sql = `EXEC dbo.GetWorkloadStatistics ?,?`
	db := clientDb.DB.Raw(sql, q.StartTime, q.EndTime).Find(&d)
	if db.Error != nil {
		return nil, fmt.Errorf("%w", db.Error)
	}
	return &d, nil
}
