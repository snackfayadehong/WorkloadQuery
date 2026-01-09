package model

type RefundNo struct {
	Yddh             string `json:"yddh" gorm:"column:yddh"`
	Rkfs             string `json:"rkfs" gorm:"column:rkfs"`
	Sczt             string `json:"sczt" gorm:"column:sczt;default:''"`
	Scsm             string `json:"scsm" gorm:"column:scsm;default:''"`
	StoreHouseName   string `json:"storeHouseName" gorm:"column:storeHouseName"`
	LeaderDepartName string `json:"leaderDepartName" gorm:"column:leaderDepartName"`
}

type RefundFullSerializer struct {
	*RefundNo
}

func (f *RefundFullSerializer) RefundSerialize() interface{} { // 返回适用于HIS接口入参的匿名结构体 用于入参序列化
	return struct {
		Yddh string `json:"yddh" gorm:"column:yddh"`
		Rkfs string `json:"rkfs" gorm:"column:rkfs"`
		Sczt string `json:"sczt" gorm:"column:sczt;default:''"`
		Scsm string `json:"scsm" gorm:"column:scsm;default:''"`
	}{
		Yddh: f.Yddh,
		Rkfs: f.Rkfs,
		Sczt: f.Sczt,
		Scsm: f.Scsm,
	}
}

type RefundSerializer interface {
	Serialize() interface{}
}

type GetRefund interface {
	GetRefundNo() error
}

type RetryRefund interface {
	RetryRefundToHis() error
}
