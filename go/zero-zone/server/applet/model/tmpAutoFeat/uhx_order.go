package tmpAutoFeat

type TmpUhxOrder struct {
	Id              int64   `json:"id" gorm:"column:id;comment:唯一id;"`
	No              string  `json:"no" label:"撤销单号" gorm:"column:no;comment:撤销单号;"`
	HxNo            string  `json:"hx_no" label:"核销单号" gorm:"column:hx_no;comment:核销单号;"`
	Typo            int     `json:"typo" label:"渠道" gorm:"column:typo;comment:渠道;"`
	CooperateShopId int64   `json:"cooperateShopId" label:"门店ID" gorm:"column:cooperate_shop_id;comment:门店ID;"`
	Content         string  `json:"content" label:"第三方返回信息" gorm:"column:content;comment:第三方返回信息;"`
	OpenNo          string  `json:"openNo" label:"平台订单" gorm:"column:open_no;comment:平台订单;"`
	Status          int     `json:"status" label:"状态" gorm:"column:status;comment:状态;"`
	IdentCode       string  `json:"identCode" label:"券码" gorm:"column:ident_code;comment:券码;"`
	QdOrderId       string  `json:"qdOrderId" label:"渠道订单号" gorm:"column:qd_order_id;comment:渠道订单号;"`
	QdGoodName      string  `json:"qdGoodName" label:"渠道商品名称" gorm:"column:qd_good_name;comment:渠道商品名称;"`
	QdPrice         float64 `json:"qdPrice" label:"渠道金额" gorm:"column:qd_price;comment:渠道金额;"`
	CreatedAt       string  `json:"createdAt" gorm:"column:created_at;comment:创建时间;"`
	UpdatedAt       string  `json:"updatedAt" gorm:"column:updated_at;comment:更新时间;"`
	DeletedAt       string  `json:"deletedAt" gorm:"column:deleted_at;comment:删除时间;"`
}
