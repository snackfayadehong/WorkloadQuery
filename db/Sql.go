package clientDb

// UserProdAcceptSql 入库信息
const UserProdAcceptSql = `select a.MEnName ,  sp.prodSpecNum as prodSpecNum, count(distinct billno) as billNum,CONVERT(varchar,CONVERT(decimal(18,2),sum(b.qty * b.BuyPrice))) as totalAmount  
from T_Prod_Enter a
INNER join T_ProdEnter_Detail b on (a.Reg_id=b.Reg_id)
 INNER  join (SELECT MEnName, SUM(SpecNum) AS prodSpecNum
FROM (
  SELECT q.MEnName ,COUNT(DISTINCT p.prod_id) AS SpecNum 
  FROM T_Prod_Enter q
  inner JOIN T_ProdEnter_Detail p   ON q.reg_id = p.reg_id
  WHERE q.billstate IN ('41', '51') AND p.IsVoid = 0 
	AND q.EnterDate >=?
	AND q.EnterDate <=?
  GROUP BY q.reg_id,q.MEnName 
) AS subQuery
GROUP BY MEnName) as sp on sp.MEnName = a.MEnName
where  a.billstate  in ('41','51')
  and b.IsVoid = 0
and a.EnterDate>= ?  and a.EnterDate<= ?
group by a.MenName , sp.prodSpecNum
`

// UserProdDpcSql 出库信息
const UserProdDpcSql = `SELECT
	a.BLMakerName,
	dp.DpSpecNum,
	COUNT ( DISTINCT a.DepartmentCollarID ) AS DpBillNum,
	CONVERT ( DECIMAL ( 18, 2 ), SUM ( b.Amount* b.RealUnitPrice ) ) AS DpTotalAmount 
FROM
	TB_DepartmentCollar a
	INNER JOIN TB_DepartmentCollarDetail b ON ( a.DepartmentCollarID= b.DepartmentCollarID ) 
	INNER JOIN (select BLMakerName, sum(specNum )  as DpSpecNum 
	from (
  SELECT q.BLMakerName ,COUNT(DISTINCT p.ProductInfoID) AS SpecNum 
  FROM TB_DepartmentCollar q
  inner JOIN TB_DepartmentCollarDetail p   ON q.DepartmentCollarID = p.DepartmentCollarID
  WHERE q.Status IN ('21', '51','61') AND p.IsVoid = 0 
	AND q.BLDate >= ?
	AND q.BLDate <= ?
	and q.TreasuryDepartment in ('200346','200418') 
  GROUP BY q.DepartmentCollarID,q.BLMakerName 
)AS subQuery
GROUP BY BLMakerName) as dp on dp.BLMakerName = a.BLMakerName
WHERE
	a.Status IN ( '21', '51', '61' ) 
  	AND b.IsVoid = 0
	AND a.TreasuryDepartment in ('200346','200418') 
	AND a.BLDate>= ?
	AND a.BLDate<= ?
GROUP BY
	a.BLMakerName,dp.DpSpecNum`

// UserProdRefundSql 退货信息
const UserProdRefundSql = `SELECT
	d.EmployeeName,
	COUNT ( b.ProductInfoID ) AS ReFSpecNum,
	COUNT ( DISTINCT a.ReturnID ) AS RefBillNum,
	CONVERT ( DECIMAL ( 18, 2 ), SUM ( b.Amount * b.UnitPrice ) ) AS RefTotalAmount 
FROM
	TB_ReturnPurchase a
	INNER JOIN TB_ReturnPurchaseDetail b ON ( a.ReturnID= b.ReturnID )
	INNER JOIN TB_Employee d ON a.BLMaker = d.HRCode 
	INNER JOIN (select BLMaker, sum(specNum )  as DpSpecNum 
	from (
  SELECT q.BLMaker,COUNT( DISTINCT p.ProductInfoID) AS SpecNum 
  FROM TB_ReturnPurchase q
  inner JOIN TB_ReturnPurchaseDetail p   ON q.ReturnID = p.ReturnID
  WHERE q.Status = '21' AND p.IsVoid = 0 
	AND q.BLDate >= ? 
	AND q.BLDate <= ?
	and q.StorehouseID  = '200346'
  GROUP BY q.ReturnID,q.BLMaker
)AS subQuery
GROUP BY BLMaker) as rt on rt.BLMaker = d.HRCode
WHERE
	a.BLDate>= ?
	AND a.BLDate< ?
	AND a.StorehouseID IN ( '200346' ) 
	AND a.Status = 21 
	AND b.IsVoid = 0 
GROUP BY
	d.EmployeeName`

// NoAccountEntrySql 科室调拨未上账单据查询
const NoAccountEntrySql = `select 
a.DepartmentCollarCode
,a.LeadingDepartmentName
,a.LeaderName
,a.BLMakerName
,CONVERT(varchar,a.CreateTime,120) as CreateTime
,CONVERT(varchar,GETDATE(),120) as StatisticalTime
,prod.Code 
,prod.Name as ProdName
,ISNULL(mode.name,'') + '|' + ISNULL(spec.Name,'') as SpecModelName
,unit.Name  as UnitName
,CONVERT(decimal(18,0),prod.ChargePrice) as ChargePrice
,CONVERT(decimal(18,0),sum(b.amount)) as Amount
from 
TB_DepartmentCollar a 
left join TB_DepartmentCollarDetail b on a.DepartmentCollarID = b.DepartmentCollarID
left join TB_ProductInfo prod on b.ProductInfoID = prod.ProductInfoID	
left join TB_SpecUnit spec on prod.Specification = spec.SpecID
left join TB_SpecUnit unit on unit.SpecID = prod.Unit
left join TB_SpecUnit mode on mode.SpecID = prod.Model
where 
a.BLDate>=? and a.BLDate <= ?
and a.Status = 61 and b.IsVoid = 0 
and TreasuryDepartment in('200416','200346','200418','200420')
group by a.DepartmentCollarCode,a.LeadingDepartmentName,a.LeaderName,a.BLMakerName,
a.CreateTime,prod.Code,prod.Name,mode.name,spec.Name
,unit.Name,prod.ChargePrice	
order by a.DepartmentCollarCode
`

// UnCheckBillSql 计费未核对数据查询
const UnCheckBillSql = `select
a.BillNo,
a.PatName,
a.[Section],
SUBSTRING(a.CgDoctor, 1, CHARINDEX('_', a.CgDoctor) - 1) AS Doctor,
a.OperateName,
a.CheckDate,
CASE 
	WHEN  DATEDIFF(day, a.CheckDate, GETDATE()) > 30 THEN '未核对(超过有效核对日期)'
	ELSE '未核对'
END AS CheckStatus
 from T_Instrument_Use a 
where a.Flag = 21 and a.Type = 2
and a.CheckDate >= ? 
and a.CheckDate <= ?
ORDER BY a.SectionId,a.CheckDate`

// NotDeliveredPurchaseSummarySql 采购订单未到货统计
const NotDeliveredPurchaseSummarySql = `select 
 ps.PurchaseSummaryID
,ps.PurchaseSummaryCode
,ps.DepartmentName
,ps.StatusName 
,ps.MakeName
,ps.AuditorDate
,ps.Remark
,CONVERT(decimal(18,2),ps.AllMoney)  as AllMoney
,ps.GoodStatusName
,ps.SupplierName
from dbo.V_PurchaseSummary ps with(nolock) 
where ps.Status in ('71','91') -- 91已下单,71部分到货
and ps.AuditorDate >= ?
and ps.AuditorDate <= ?
and ps.Remark not like '%高低储%'
order by ps.SupplierName`

// NotDeliveredPurchaseSummaryDetailSql 采购订单未到货明细统计
const NotDeliveredPurchaseSummaryDetailSql = `select  
 psd.PurchaseSummaryID
,prod.Code
,prod.Name as 'ProdName'
,prod.HospitalSpec
,ent.EnterpriseName
,CONVERT(decimal(18,2),psd.UnitPrice) as 'UnitPrice'
,spec.Name as 'SpecName'
,CONVERT(decimal(18,2),psd.Amount) as 'Amount'
,CONVERT(decimal(18,2),psd.FactInAmount) as 'FactInAmount'
,CONVERT(decimal(18,2),psd.RefundAmount)  as 'RefundAmount'
,CONVERT(decimal(18,2),psd.Amount - psd.FactInAmount) as 'NotDeliveredAmount'
,psd.Remark
from dbo.TB_PurchaseSummaryDetail psd with(nolock)
left join TB_ProductInfo prod on prod.ProductInfoID = psd.ProductInfoID
left join TB_EnterpriseInfo ent on ent.EnterpriseID = prod.DefaultSupplierID
left join TB_SpecUnit spec on spec.SpecID = psd.Unit
where psd.IsVoid = 0 -- 明细有效
and psd.Amount != psd.FactInAmount -- 未到货
and PurchaseSummaryID in (?)`
