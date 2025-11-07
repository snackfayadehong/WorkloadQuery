package model

type DeliveryNo struct {
	Ckdh       string `json:"yddh" gorm:"column:ckdh"`
	DetailSort string `json:"-" gorm:"column:detailSort"`
	Ckfs       string `json:"ckfs" gorm:"column:ckfs"`
	Sczt       string `json:"sczt" gorm:"column:sczt;default:''"`
	Scsm       string `json:"scsm" gorm:"column:scsm;default:''"`
}

type Get interface {
	GetDeliveryNo() error
}

type ReTry interface {
	DeliveryNoRetryToHis() error
}
