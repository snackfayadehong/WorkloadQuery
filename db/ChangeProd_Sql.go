package clientDb

/*
	基础字典修改接口
*/

// QueryProd 查询产品信息
const QueryProd = `SELECT
    prod.ProductInfoID -- 产品ID
	,prod.Code -- 院内代码
	,prod.HisProductCode3 AS HospitalName -- 院内名称
	,prod.HospitalSpec -- 院内规格
	,prod.HisProductCode7 AS YGCGID -- 网采平台产品ID
	,prod.HisProductCode7Source -- 集采平台产品临时ID 审核过后会改变HisProductCode7值
	,prod.HisProductCode7Status -- 集采平台产品ID状态
	,prod.CusCategoryCode -- 104分类编码
	,
	CASE
	WHEN cus.SubjectGrade = 0 THEN '' 
	WHEN cus.SubjectGrade = 1 THEN '00' 
	ELSE SUBSTRING ( prod.CusCategoryCode, 1, ( CONVERT ( INT, cus.SubjectGrade ) - 1 ) * 2 ) + '0000'
	END ParentCusCategoryCode
	,td.TenderCode AS TradeCode -- 交易编码
    --,case 
    --when yb.MedicareCode = '' then yb.ChargeRuleID
    --when gjyb.MedicareCode = '' then gjyb.ChargeRuleID
    --else ''
       -- end	as MedicareID 
	--,yb.MedicareCode -- 医保编码
	--,gjyb.MedicareCode as CountryMedicareCode -- 国家医保
    ,prod.ChargePrice -- 收费价格
	,jc.SysCode -- 集采系统编码
	,jc.SysId -- 集采系统编号
	,prod.IsVoid -- 0：启用 1：停用
	,ep.PurState -- 0：供货 1：停止供货
FROM
	TB_ProductInfo  prod WITH (NOLOCK)
	LEFT JOIN TB_TenderCode td  WITH (nolock) ON td.ProductInfoID = prod.ProductInfoID 
	AND td.IsVoid = 0 
	AND td.MedicareType = 1
	-- 医保编码
	LEFT JOIN TB_ProductChargeRule yb WITH (nolock) ON yb.ProductInfoID = prod.ProductInfoID 
	AND yb.IsVoid = 0 
	and yb.MedicareType = 1 
	-- 国家医保编码
	LEFT JOIN TB_ProductChargeRule gjyb WITH (nolock) on yb.ProductInfoID = prod.ProductInfoID
	and  yb.IsVoid = 0
	and yb.MedicareType = 3
	LEFT JOIN TB_EnterpriseProduct ep  WITH (nolock) ON ep.ProductID = prod.ProductInfoID 
	AND prod.DefaultSupplierID = ep.EnterpriseID
	LEFT JOIN TB_ProductInfoJCSysCode jc WITH (nolock) ON jc.Prod_Id = prod.ProductInfoID 
	LEFT JOIN TB_ProductCustomCategory cus WITH (nolock) on prod.CusCategoryCode = cus.CusCategoryCode
WHERE
	prod.CODE IN (?)`
