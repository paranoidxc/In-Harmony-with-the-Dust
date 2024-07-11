package tmpAutoFeat

import (
	"time"

	"gorm.io/gorm"
)

//gorm生成表结构时 结构体不要加 tmp  并且结构体与文件名必须对应，文件名要用下划线 如：test_gorm  ID必须固定为ID 类型固定int64，否则生成之后还要稍微修改

type TestGorm struct {
	ID        int64          `json:"id" gorm:"primarykey;comment:唯一id;"`                                          // 主键ID
	CreatedAt time.Time      `json:"created_at" gorm:"comment:创建时间;" time_format:"sql_datetime" time_utc:"false"` // 创建时间
	UpdatedAt time.Time      `json:"updated_at" gorm:"comment:更新时间;" time_format:"sql_datetime" time_utc:"false"` // 更新时间
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index;comment:删除时间;"`                                       // 删除时间
	Text      string         `json:"text" gorm:"type:varchar(255);comment:文本;"`                                   // 文本
}
