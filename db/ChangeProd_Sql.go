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

const ProductInfo_UpdatePostDataSQL = `SELECT  
a.code AS ypdm,  
REPLACE(a.Name, CHAR(13) + CHAR(10), ' ') as ypmc,
isnull( d.Name, '' ) AS yppp,  
LEFT( isnull( gi.PinYin, gi2.PinYin ), 11) AS zjm,  
LEFT(ISNULL(gi2.Name, a.Name), 11) as ypbm,
LEFT(ISNULL(gi2.PinYin, gi.PinYin), 11) as pym,  
left(isnull(gi2.name,a.name	),11) as ypbm1,
left(isnull(gi2.PinYin, gi.PinYin), 11) as pym1,
left(isnull(gi2.name,a.name	),11) as ypbm2,
left(isnull(gi2.PinYin, gi.PinYin), 11) as pym2,
case 
when ISNULL(b.Name,'') = ''and ISNULL(mo.Name,'') = '' then a.HospitalSpec
when ISNULL(mo.Name,'') != '' and ISNULL(b.name,'') != '' then mo.Name + '|' + b.Name
when ISNULL(b.Name,'') = '' then mo.Name
else b.name
end ypgg,
CONVERT(VARCHAR(20),cate.Remark) AS yplb,  
'17' jxbm,  
'09' AS lbdm,
a.HisProductStoreCode AS kfdm,  
a.TopLimit AS kfgx,  
a.DownLimit AS kfdx, 
Cr.MedicareCode ypbwm,
CASE  
WHEN charindex( '无注册证', isnull( lmr.RegistNum, '' ) ) > 0 THEN  'false'   
WHEN charindex( '进', isnull( lmr.RegistNum, '' ) ) > 0 THEN  'true'   -- 进口  
WHEN charindex( '许', isnull( lmr.RegistNum, '' ) ) > 0 THEN  'false' ELSE 'false'  -- 港澳台  国产   
END sfwjkcl, 
tcode.TenderCode lsh,
a.CommonName tymc,
left(lmr.RegistNum,15) ypzczh,
CONVERT(varchar ,lmr.EndDate,23) AS ypzczhxq,
'' gnzdl,
CASE   
WHEN a.ChargeStatus is NULL THEN '0'  
WHEN a.ChargeStatus = '0' THEN '1'  
WHEN a.ChargeStatus = '1' THEN '0'  
END as cljflx,
case 
when a.Buy = 0 then  1 
when a.Buy = 1 then  2
ELSE 3 
END sfwwhp1,
a.PackingInstruction ypbz,
a.StorageCondition cctj,
CAST ( a.ProductInfoID AS VARCHAR ( 100 ) ) + CAST ( En1.EnterpriseID AS VARCHAR ( 100 ) ) AS ypspdm,
CONVERT(varchar ,a.UpdateTime,120) AS lrrq,
 --a.IsVoid sybz,
case when a.isvoid = 0 then '1'
when a.IsVoid = 1 then '0'
else '1'
end as 'sybz',
tcode.TenderCode  AS wcspdm,  
a.HisProductCode7Source AS wccpid,  
'' AS kfbz,
a.OpenTender AS sfjc,
a.isJcSelect AS jcsfzb,
cus.CusCategoryName as clfl104,
jcxx.SysCode as jcptxtbh,
jcxx.SysId as jcptxtbm,
'' as xgczy,
 isnull( a.ForbinNotBak, '' ) bz   
FROM  
 TB_ProductInfo a  
 LEFT JOIN TB_SpecUnit mo on a.Model = mo.SpecID ----- 型号
 LEFT JOIN TB_SpecUnit b ON a.Specification= b.SpecID -----规格  
 LEFT JOIN TB_SpecUnit tu ON a.MatInUnit= tu.SpecID ------ 仓库存放单位  
 LEFT JOIN TB_ProductCustomCategory c ON a.CategoryCode= c.CusCategoryCode  
 LEFT JOIN TB_BrandInfo d ON a.BrandID= d.BrandID -----产品品牌名称  
 LEFT JOIN TB_EnterpriseProduct ep ON a.ProductInfoID= ep.ProductID   
 AND ep.EnterpriseID= a.DefaultSupplierID ---供货关系  
 --LEFT JOIN TB_EnterpriseInfo En ON En.Type= 1   
 --AND En.EnterpriseID = a.EnterpriseID --------生产企业  
 LEFT JOIN TB_EnterpriseInfo En1 ON En1.EnterpriseID = a.DefaultSupplierID --------默认供应商  
 LEFT Join TB_ProductChargeRule Cr on cr.ProductInfoID = a.ProductInfoID  and Cr.IsVoid = 0 and Cr.MedicareType = 1 -- 医保  
 LEFT JOIN TB_GeneralInfo gi ON gi.IsDefault= '1'  -- 产品名称别名  
 and  a.ProductInfoID= gi.ProductInfoID  
 LEFT JOIN (select * from (  
select *,row_number() over(partition by productinfoid order by generalid asc) rn from TB_GeneralInfo  
where IsDefault != '1' and isvoid = 0  
) tn where tn.rn=1  
 )gi2 on gi2.ProductInfoID = a.ProductInfoID  
 LEFT JOIN TB_TenderCode tcode ON tcode.ProductInfoID= a.ProductInfoID ----中标  
 LEFT JOIN (  
 SELECT  
  rp.ProductID,  
  lm.LicenseD,  
  lm.RegistNum,  
  lm.IsVoid,  
  lm.EndDate,  
  rn = ROW_NUMBER ( ) OVER ( PARTITION BY rp.ProductID ORDER BY lm.EndDate DESC )   
 FROM  
  TB_RegisterProduct rp  
  INNER JOIN TB_LicenceMerge lm ON lm.LicenceType = 6   
  AND lm.LicenseD = rp.LicenseD   
 WHERE  
  rp.IsVoid = '0'   
 ) lmr ON lmr.rn = 1   
 AND lmr.ProductID = a.ProductInfoID ------最新注册证  
left join TB_ProductCategory cate on cate.CategoryCode = a.CategoryCode   
and cate.SubjectGrade = 3   
and cate.IsVoid = 0
left join TB_ProductCustomCategory cus on cus.CusCategoryCode = a.CusCategoryCode
left join TB_ProductInfoJCSysCode jcxx on jcxx.prod_id = a.ProductInfoID  and jcxx.IsVoid = 0
where a.Code = ?`
