package tmpAutoFeat

// gorm:"column:firm_name;not null;uniqueIndex;type:varchar(255);comment:厂商名称;"`
// gorm tag 暂时只用到 column 和 comment

type TmpDemoCurd struct {
	Id int64 `json:"id" gorm:"column:id;comment:唯一id;"`

	FirmName  string `json:"firmName" validate:"required" label:"电影名称" gorm:"column:firm_name;comment:厂商名称;"`
	FirmAlias string `json:"firmAlias" gorm:"column:firm_alias;comment:厂商别名;"`
	FirmCode  string `json:"firmCode" validate:"number,gte=1" label:"电影编码" gorm:"column:firm_code;comment:电影编码;"`
	FirmDesc  string `json:"firmDesc" gorm:"column:firm_desc;comment:厂商描述;"`

	CreatedAt string `json:"createdAt" gorm:"column:create_at;comment:创建时间;"`
	UpdatedAt string `json:"updatedAt" gorm:"column:update_at;comment:更新时间;"`
	DeletedAt string `json:"deletedAt" gorm:"column:delete_at;comment:删除时间;"`
}
