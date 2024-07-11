func (m *default{{.upperStartCamelObject}}Model) Update(data *{{.upperStartCamelObject}}) error {
	result := m.db.Model(&{{.upperStartCamelObject}}{}).Where("{{.originalPrimaryKey}} = ?", data.ID).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
