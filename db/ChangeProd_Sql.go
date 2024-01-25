package clientDb

/*
	基础字典修改接口
*/

// QueryProd 查询产品信息
const QueryProd = `SELECT
	prod.CODE -- 院内代码
	,prod.HisProductCode3 AS HospitalName -- 院内名称
	,prod.HospitalSpec -- 院内规格
	,prod.HisProductCode7 AS YGCGID -- 网采平台产品ID
	,prod.HisProductCode7Source -- 集采平台产品临时ID 审核过后会改变HisProductCode7值
	,prod.HisProductCode7Status -- 集采平台产品ID状态
	,prod.CusCategoryCode -- 104分类编码
	,td.TenderCode AS TradeCode -- 交易编码
	,yb.MedicareCode -- 医保编码
	,jc.SysCode -- 集采系统编码
	,jc.SysId -- 集采系统编号
	,prod.IsVoid -- 0：启用 1：停用
	,ep.PurState -- 0：供货 1：停止供货
FROM
	TB_ProductInfo prod
	LEFT JOIN TB_TenderCode td ON td.ProductInfoID = prod.ProductInfoID 
	AND td.IsVoid = 0 
	AND td.MedicareType = 1
	LEFT JOIN TB_ProductChargeRule yb ON yb.ProductInfoID = prod.ProductInfoID 
	AND td.IsVoid = 0 
	AND yb.MedicareType = 1
	LEFT JOIN TB_EnterpriseProduct ep ON ep.ProductID = prod.ProductInfoID 
	AND prod.DefaultSupplierID = ep.EnterpriseID
	LEFT JOIN TB_ProductInfoJCSysCode jc ON jc.Prod_Id = prod.ProductInfoID 
WHERE
	prod.CODE IN (?)`
