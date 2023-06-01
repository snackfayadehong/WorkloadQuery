package clientDb

import (
	"fmt"
)

type UserProdAccept struct {
	Oper   string // 名称
	Amount int
	// ProdAcSpec   int    // 品规数
	// ProdAcBillNo int    // 单据数
	// ProdAcTotal  int64  // 总金额
}

// ProdAcceptQuery { 入库信息查询
func ProdAcceptQuery(startTime string, endTime string) {
	var err error
	fmt.Println(err)
	// var rows *sql.Rows
	var ProdAccept UserProdAccept
	// 品规数
	DB.Raw(`select oper, num  from (
select a.Oper ,COUNT( b.SpecModelName )as Num from T_Prod_Enter a
left join T_ProdEnter_Detail b on (a.Reg_id=b.Reg_id)
 where  a.billstate  in ('41','51') and a.oper ='廖小凤'
and a.EnterDate>= ? and a.EnterDate<= ?
group by a.Oper )b
`, startTime, endTime).Scan(&ProdAccept)
	// if err != nil {
	// 	return
	// }
	// rows.Scan(&prodAccept)
	fmt.Println(ProdAccept)

	// 单据数
	// 	rows, err = DB.Raw(`select oper ,count(billno) from (
	// select a.Oper, a.billno from T_Prod_Enter  a  where  a.billstate in('41','51')
	// and a.EnterDate>= ? and a.EnterDate<= ?
	// group by a.Oper,a.billno ) b
	// group by Oper`, startTime, endTime).Rows()
	//
	// 	// 总金额
	// 	rows, err = DB.Raw(`select a.Oper,sum( b.Qty * b.BuyPrice ) as totalPrice from T_Prod_Enter  a
	// left Join T_ProdEnter_Detail b on a.Reg_id = b.Reg_id
	//  and a.billstate in('41','51')
	// and a.EnterDate>=? and a.EnterDate<?
	// where b.Qty != 0
	// group by a.oper `, startTime, endTime).Rows()
	// 	defer func(rows *sql.Rows) {
	// 		err := rows.Close()
	// 		if err != nil {
	// 			return
	// 		}
	// 	}(rows)

}
