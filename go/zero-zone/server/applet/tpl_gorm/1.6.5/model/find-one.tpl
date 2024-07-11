func (m *default{{.upperStartCamelObject}}Model) FindOne({{.lowerStartCamelPrimaryKey}} int64) (*{{.upperStartCamelObject}}, error) {
	var resp {{.upperStartCamelObject}}
	// 使用 GORM 根据ID查询记录
	result := m.db.First(&resp, {{.lowerStartCamelPrimaryKey}})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		} else {
			return nil, result.Error
		}
	}
	return &resp, nil
}

func (m *default{{.upperStartCamelObject}}Model) FindAllByWhere(where string) ([]*{{.upperStartCamelObject}}, error) {
	var resp []*{{.upperStartCamelObject}}
	// 使用 GORM 执行查询
	result := m.db.Where(where).Order("{{.originalPrimaryKey}} DESC").Find(&resp)
	if result.Error != nil {
		return nil, result.Error
	}
	return resp, nil
}

func (m *default{{.upperStartCamelObject}}Model) FindAllByWhereCount(where string) (int64, error) {
	var resp int64
	// 使用 GORM 执行查询并返回计数
	result := m.db.Model(&{{.upperStartCamelObject}}{}).Where(where).Count(&resp)
	if result.Error != nil {
		return 0, result.Error
	}
	return resp, nil
}

func (m *default{{.upperStartCamelObject}}Model) FindPageByWhere(where string, page int64, limit int64) ([]*{{.upperStartCamelObject}}, error) {
	var resp []*{{.upperStartCamelObject}}
	// 使用 GORM 执行分页查询
	result := m.db.Where(where).Offset(int((page - 1) * limit)).Limit(int(limit)).Order("{{.originalPrimaryKey}} DESC").Find(&resp)
	if result.Error != nil {
		return nil, result.Error
	}
	return resp, nil	
}

func (m *default{{.upperStartCamelObject}}Model) FindPageByWhereCount(where string) (int64, error) {
	var resp int64
	// 使用 GORM 执行查询并返回计数
	result := m.db.Model(&{{.upperStartCamelObject}}{}).Where(where).Count(&resp)
	if result.Error != nil {
		return 0, result.Error
	}
	return resp, nil
}