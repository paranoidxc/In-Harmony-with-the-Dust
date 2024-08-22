package tmpAutoFeat

type TmpSysRegion struct {
	Id        int64  `json:"id" db:"column:id;comment:唯一id;"`
	No        string `json:"no" label:"编号" db:"column:no;comment:编号;"`
	Name      string `json:"name" label:"名称" db:"column:no;comment:名称;"`
	ParentNo  string `json:"parentNo" label:"父级编号" db:"column:no;comment:父级编号;"`
	Code      string `json:"code" label:"区码" db:"column:no;comment:区码;"`
	TypeName  string `json:"typeName" label:"类型名称" db:"column:no;comment:类型名称;"`
	PySzm     string `json:"pYSzm" label:"拼音" db:"column:no;comment:拼音""`
	CreatedAt string `json:"createdAt" db:"column:created_at;comment:创建时间;"`
	UpdatedAt string `json:"updatedAt" db:"column:updated_at;comment:更新时间;"`
	DeletedAt string `json:"deletedAt" db:"column:deleted_at;comment:删除时间;"`
	IsDel     int64  `json:"isDel" db:"column:deleted_at;comment:删除时间;"`
}
