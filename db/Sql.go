package clientDb

// UserProdAcceptSql 入库信息
const UserProdAcceptSql = `
select a.Oper ,COUNT( b.SpecModelName )as prodSpecNum ,count(distinct billno) as billNum,CONVERT(varchar,CONVERT(decimal(18,4),sum(b.qty * b.BuyPrice))) as totalAmount from T_Prod_Enter a
left join T_ProdEnter_Detail b on (a.Reg_id=b.Reg_id)
where  a.billstate  in ('41','51')
and a.EnterDate>= ? and a.EnterDate<= ?
group by a.Oper`

// UserProdDpcSql 出库信息
const UserProdDpcSql = `select a.BLMakerName,COUNT(b.ProductInfoID) as DpSpecNum ,
count( distinct a.DepartmentCollarID) as  DpBillNum, CONVERT(decimal(18,4),sum(b.Amount*b.RealUnitPrice))as DpTotalAmount from TB_DepartmentCollar a
left join TB_DepartmentCollarDetail b on (a.DepartmentCollarID=b.DepartmentCollarID)
left join TB_ProductInfo c on (b.ProductInfoID=c.ProductInfoID)
where a.Status in ('21','51','61')  and a.TreasuryDepartment = '200346'
and  a.BLDate>= ? and a.BLDate<=?
group by a.BLMakerName
`

// UserProdRefundSql 退货信息
const UserProdRefundSql = `select d.EmployeeName, count(b.ProductInfoID) as ReFSpecNum, count(distinct a.ReturnID) as RefBillNum,
CONVERT(decimal(18,4),sum(b.Amount * b.UnitPrice)) as RefTotalAmount
 from TB_ReturnPurchase a
left join TB_ReturnPurchaseDetail b on (a.ReturnID=b.ReturnID)
left join TB_ProductInfo c on (b.productinfoid=c.productinfoid)
left join TB_Employee d on a.BLMaker = d.HRCode
where a.BLDate>=? 
and a.BLDate<? 
and a.StorehouseID in ('200346')
and a.Status = 21
and b.IsVoid = 0
group by d.EmployeeName`
