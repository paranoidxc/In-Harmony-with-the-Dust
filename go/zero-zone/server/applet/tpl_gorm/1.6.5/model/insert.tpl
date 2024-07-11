func (m *default{{.upperStartCamelObject}}Model) Insert(data *{{.upperStartCamelObject}}) (error) {
	result := m.db.Create(data)
	if result.Error != nil {
		return result.Error
	}
	return nil	
}
