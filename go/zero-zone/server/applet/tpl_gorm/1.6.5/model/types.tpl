type (
	{{.lowerStartCamelObject}}Model interface{
		Insert(data *{{.upperStartCamelObject}}) error
		FindOne(id int64) (*{{.upperStartCamelObject}}, error)
		Update(data *{{.upperStartCamelObject}}) error
		Delete(id int64) error
		Deletes(ids []int64) error
		FindAllByWhere(where string) ([]*{{.upperStartCamelObject}}, error)
		FindAllByWhereCount(where string) (int64, error)
		FindPageByWhere(where string, page int64, limit int64) ([]*{{.upperStartCamelObject}}, error)
		FindPageByWhereCount(where string) (int64, error)			}

	default{{.upperStartCamelObject}}Model struct {
		db    *gorm.DB
		table string
	}

	{{.upperStartCamelObject}} struct {
		replace string
	}
)
