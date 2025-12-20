package controller

import (
	clientDb "WorkloadQuery/db"
	"WorkloadQuery/model"
	"fmt"
	"strconv"
)

type DictCompareController struct{}

func (ctrl *DictCompareController) GetLocalDictInfo(keyword string) ([]model.LocalDictRow, error) {
	var rows []model.LocalDictRow // 修改点 2：变量名复数化

	// ypgg 的 CASE 逻辑片段
	ypggSql := `CASE 
      WHEN ISNULL(SP.Name,'') = '' AND ISNULL(mo.Name,'') = '' THEN a.HospitalSpec
      WHEN ISNULL(mo.Name,'') != '' AND ISNULL(SP.name,'') != '' AND ISNULL(mo.Name,'') != ISNULL(SP.name,'') THEN mo.Name + '|' + SP.Name
      WHEN ISNULL(mo.Name,'') != '' AND ISNULL(SP.name,'') != '' AND ISNULL(mo.Name,'') = ISNULL(SP.Name,'') THEN mo.Name
      WHEN ISNULL(SP.Name,'') = '' THEN mo.Name
      ELSE SP.name
    END AS ypgg`

	// 链式调用
	query := clientDb.DB.Table("TB_ProductInfo AS A").
		Select("A.Code AS ypdm, A.ProductInfoID, A.Name AS ypmc, " + ypggSql + ", SPEC.Name AS kfdw, A.PurchasePrice AS kfcgj, A.ChargePrice AS kflsj, A.HisProductStoreCode AS kfdm, B.SourceValue AS gsdm,  B.EnterpriseName AS ghdw").
		Joins("JOIN TB_EnterpriseInfo B ON B.EnterpriseID = A.DefaultSupplierID").
		Joins("LEFT JOIN TB_SpecUnit SPEC ON SPEC.SpecID = A.UseUnit").
		Joins("LEFT JOIN TB_SpecUnit MO ON MO.SpecID = A.Model").
		Joins("LEFT JOIN TB_SpecUnit SP ON SP.SpecID = A.Specification")

	// 1. 尝试将关键字转化为数字
	_, err := strconv.Atoi(keyword)
	if err == nil {
		// 根据keyword长度 判断查询 ProductInfoID 还是Code
		if len(keyword) <= 6 {
			query = query.Where("A.ProductInfoID = ?", keyword)
		} else {
			query = query.Where("A.Code = ?", keyword)
		}
	} else {
		return nil, fmt.Errorf("非法字符")
	}
	err = query.Order("A.ProductInfoID").Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}
