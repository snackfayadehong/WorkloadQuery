package clientDb

// QueryBillNo 领用出库未推送数据
const QueryBillNo = `select 
	'01' ckfs ,
	 delivery.DeliveryID as ckdh,
	deliveryRecord.DetailSort as detailSort
	from TB_DeliveryApplyDetailRecord deliveryRecord with(nolock) 
	join TB_DeliveryApply delivery with(nolock) on deliveryRecord.DeliveryID = delivery.DeliveryID
	where 1=1
	and deliveryRecord.IsVoid = 0
	and delivery.Source = '1' 
	AND IsStockGoods = '0' 
	AND Delivery.[Type] = '1' 
	AND Delivery.[Status] IN (61, 71, 41, 81, 22, 91, 19, 29, 99)  
	AND (Delivery.[IsStockGoods] <> '1' OR Delivery.[IsStockGoods] IS NULL)
	AND ISNULL(deliveryRecord.OutNumber,'') = ''
	AND deliveryRecord.CreateTime >= ? --AND deliveryRecord.CreateTime <= ?
	group by 
	deliveryRecord.DetailSort,
	delivery.DeliveryID 
 `

// UpdateDelivery_Sql 领用出库推送后修改状态
const UpdateDelivery_Sql = `UPDATE TB_DeliveryApplyDetailRecord set OutNumber = ?  WHERE DeliveryID = ? And DetailSort = ?`
