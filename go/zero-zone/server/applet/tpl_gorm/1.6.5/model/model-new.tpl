func new{{.upperStartCamelObject}}Model(gorm *gorm.DB) *default{{.upperStartCamelObject}}Model {
	return &default{{.upperStartCamelObject}}Model{
		table:      {{.table}},
		db:    gorm,
	}
}

