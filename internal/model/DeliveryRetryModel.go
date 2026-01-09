package model

type DeliveryNo struct {
	Ckdh             string `json:"yddh" gorm:"column:ckdh"`
	DetailSort       string `json:"detailSort" gorm:"column:detailSort"` // 出库明细序号
	Ckfs             string `json:"ckfs" gorm:"column:ckfs"`
	Sczt             string `json:"sczt" gorm:"column:sczt;default:''"`
	Scsm             string `json:"scsm" gorm:"column:scsm;default:''"`
	DeliveryCode     string `json:"deliveryCode" gorm:"column:deliveryCode"`         // 出库单号
	StoreHouseName   string `json:"storeHouseName" gorm:"column:storeHouseName"`     // 供货库房
	LeaderDepartName string `json:"leaderDepartName" gorm:"column:leaderDepartName"` // 领用二级库房
}
type DeliverSerializer interface {
	Serialize() interface{}
}

type DeliveryFullSerializer struct {
	*DeliveryNo
}

func (f *DeliveryFullSerializer) DeliverySerialize() interface{} { // 返回适用于HIS接口入参的匿名结构体 用于入参序列化
	return struct {
		Ckdh string `json:"yddh" gorm:"column:ckdh"`
		Ckfs string `json:"ckfs" gorm:"column:ckfs"`
		Sczt string `json:"sczt" gorm:"column:sczt;default:''"`
		Scsm string `json:"scsm" gorm:"column:scsm;default:''"`
	}{
		Ckdh: f.Ckdh + f.DetailSort,
		Ckfs: f.Ckfs,
		Sczt: f.Sczt,
		Scsm: f.Scsm,
	}
}

type Get interface {
	GetDeliveryNo() error
}

type ReTry interface {
	DeliveryNoRetryToHis() error
}
