package model

// ChangeInfoElement 字典信息入参
type ChangeInfoElement struct {
	Code             string `json:"Code,required"`             // 院内代码
	HospitalName     string `json:"HospitalName,omitempty"`    // 院内产品名称
	HospitalSpec     string `json:"HospitalSpec,omitempty"`    // 院内规格
	YGCGID           string `json:"YGCGID,omitempty"`          // 网采平台产品ID
	TradeCode        string `json:"TradeCode,omitempty"`       // 商品代码，商品代码或挂网流水号（平台供货商填写交易编码的内容）
	MedicareCode     string `json:"MedicareCode,omitempty"`    // 医保代码
	EighteenProdType string `json:"EighteenProdType,required"` // 18类重点监控序号
	CategoryCode     string `json:"CategoryCode,omitempty"`    // 18类分类代码，匹配104分类三级目录,如果匹配成功不修改，否则修改为三级目录
	SysCode          string `json:"SysCode,omitempty"`         // 系统编码
	SysID            string `json:"SysID,omitempty"`           // 系统编号
	OpenTender       string `json:"OpenTender,omitempty"`      // 集采状态  1集采  0非集采
	SupplyStatus     string `json:"SupplyStatus,omitempty"`    // 供货状态
	HRCode           string `json:"HRCode,required"`           // 修改人员工号
	Remark           string `json:"Remark,omitempty"`          // 备注
}

// ExceptionProd 异常字典记录
type ExceptionProd struct {
	Code     string
	PurState int
	IsVoid   int
	HasError bool
}

// ProdMedicareCode 医保代码
type ProdMedicareCode struct {
	ProductInfoID    int    `gorm:"column:ProductInfoID"`
	ChargeRuleID     int    `gorm:"column:ChargeRuleID"`
	MedicareCode     string `gorm:"column:MedicareCode"`
	MedicareCodeTemp string `gorm:"column:MedicareCode"`
}

// ProductInfo 查询数据库产品字典信息
type ProductInfo struct {
	ProductInfoID         int    `gorm:"column:ProductInfoID"`         // 产品ID
	Code                  string `gorm:"column:Code"`                  // 院内代码
	HospitalName          string `gorm:"column:HospitalName"`          // 院内名称
	HospitalSpec          string `gorm:"column:HospitalSpec"`          // 院内规格
	OpenTender            string `gorm:"column:OpenTender"`            // 集采状态
	YGCGID                string `gorm:"column:YGCGID"`                // 网采平台产品ID
	HisProductCode7Source string `gorm:"column:HisProductCode7Source"` // 集采平台产品临时ID 审核过后会改变HisProductCode7值
	HisProductCode7Status int    `gorm:"column:HisProductCode7Status"` // 集采平台产品ID状态
	CusCategoryCode       string `gorm:"column:CusCategoryCode"`       // 104分类编码
	ParentCusCategoryCode string `gorm:"column:ParentCusCategoryCode"` // 104分类编码(第3级)
	TradeCode             string `gorm:"column:TradeCode"`             // 交易编码
	ChargePrice           string `gorm:"column:ChargePrice"`           // 收费价格
	SysCode               string `gorm:"column:SysCode"`               // 集采系统编码
	SysId                 string `gorm:"column:SysId"`                 // 集采系统编号
	IsVoid                int    `gorm:"column:IsVoid"`                // 0：启用 1：停用
	PurState              int    `gorm:"column:PurState"`              // 0：供货 1：停止供货
}

type GetProduct interface {
	GetProductInfo(Where []string) ([]ProductInfo, error)
}
type ChangeProduct interface {
	ChangeMisProductInfo([]ProductInfo, string) error
	ChangeHisProductInfo(Code string, HrCode string) error
}

// ChangeHisProductInfoModel MIS-proc-Edit-cljjxx 柯林布瑞业务中台 HIS材料信息修改接口
type ChangeHisProductInfoModel struct {
	Ypdm     string `json:"ypdm" gorm:"column:ypdm"` // 院内代码
	Ypmc     string `json:"ypmc" gorm:"column:ypmc"` // 药品名称
	Yppp     string `json:"yppp" gorm:"column:yppp"` // 药品品牌
	Zjm      string `json:"zjm" gorm:"column:zjm"`   // 助记码
	Ypbm     string `json:"ypbm" gorm:"column:ypbm"` // 药品别名
	Pym      string `json:"pym" gorm:"column:pym"`   // 拼音码
	Ypbm1    string `json:"ypbm1" gorm:"column:ypbm1"`
	Pym1     string `json:"pym1" gorm:"column:pym1"`
	Ypbm2    string `json:"ypbm2" gorm:"column:ypbm2"`
	Pym2     string `json:"pym2" gorm:"column:pym2"`
	Ypgg     string `json:"ypgg" gorm:"column:ypgg"`         // 药品规格
	Yplb     string `json:"yplb" gorm:"column:yplb"`         // 药品类别 030106
	Jxbm     string `json:"jxbm" gorm:"column:jxbm"`         // 药品剂型（试剂、器械；默认为“17”其他	）
	Lbdm     string `json:"lbdm" gorm:"column:lbdm"`         // 类别（卫生材料、其他材料；默认为“09”卫生）
	Kfdm     string `json:"kfdm" gorm:"column:kfdm"`         // 库房代码
	Kfgx     int    `json:"kfgx" gorm:"column:kfgx"`         // 库房高限
	Kfdx     int    `json:"kfdx" gorm:"column:kfdx"`         // 库房底限
	Ypbwm    string `json:"ypbwm" gorm:"column:ypbwm"`       // 医保对码
	Sfwjkcl  string `json:"sfwjkcl" gorm:"column:sfwjkcl"`   // 是否为进口材料
	Lsh      string `json:"lsh" gorm:"column:lsh"`           // 流水号
	Tymc     string `json:"tymc" gorm:"column:tymc"`         // 材料通用名称
	Ypzczh   string `json:"ypzczh" gorm:"column:ypzczh"`     // 注册证号
	Ypzczhxq string `json:"ypzczhxq" gorm:"column:ypzczhxq"` // 注册证效期
	Gnzdl    string `json:"gnzdl" gorm:"column:gnzdl"`       // 国内总代
	Cljflx   string `json:"cljflx" gorm:"column:cljflx"`     // 材料计费类型
	Sfwwhp1  string `json:"sfwwhp1" gorm:"column:sfwwhp1"`   // 审核材料类别
	Ypbz     string `json:"ypbz" gorm:"column:ypbz"`         // 材料包装
	Cctj     string `json:"cctj" gorm:"column:cctj"`         // 储存要求
	Ypspdm   string `json:"ypspdm" gorm:"column:ypspdm"`     // 怡道材料代码
	Lrrq     string `json:"lrrq" gorm:"column:lrrq"`         // 录入日期/修改日期
	Sybz     string `json:"sybz" gorm:"column:sybz"`         // 启用标志
	Wcspdm   string `json:"wcspdm" gorm:"column:wcspdm"`     // 网采商品代码
	Wccpid   string `json:"wccpid" gorm:"column:wccpid"`     // 网采产品ID
	Kfbz     string `json:"kfbz" gorm:"column:kfbz"`         // 18类国家重点监控耗材
	Sfjc     string `json:"sfjc" gorm:"column:sfjc"`         // 是否集采
	Jcsfzb   string `json:"jcsfzb" gorm:"column:jcsfzb"`     // 集采是否中标
	Clfl104  string `json:"clfl104" gorm:"column:clfl104"`   // 104 分类代码
	Jcptxtbh string `json:"jcptxtbh" gorm:"column:jcptxtbh"` // 集采平台系统编号
	Jcptxtbm string `json:"jcptxtbm" gorm:"column:jcptxtbm"` // 集采平台系统编码
	Xgczy    string `json:"xgczy" gorm:"column:xgczy"`       // 修改人员代码
	Bz       string `json:"bz" gorm:"column:bz"`             // 备注
}
