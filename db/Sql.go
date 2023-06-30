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
	and q.TreasuryDepartment = '200346'
  GROUP BY q.DepartmentCollarID,q.BLMakerName 
)AS subQuery
GROUP BY BLMakerName) as dp on dp.BLMakerName = a.BLMakerName
WHERE
	a.Status IN ( '21', '51', '61' ) 
	AND a.TreasuryDepartment = '200346' 
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
const NoAccountEntrySql = `SELECT 
DepartmentCollarCode
,BLDate
,LeadingDepartmentName
,LeaderName
,TreasuryDepartmentName
,BLMakerName
from TB_DepartmentCollar 
where 
TreasuryDepartment = '200346' and Status = 61
and BLDate <= ? 
Order by DepartmentCollarCode`
