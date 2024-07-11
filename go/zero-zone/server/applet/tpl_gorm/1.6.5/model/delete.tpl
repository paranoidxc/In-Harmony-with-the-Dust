func (m *default{{.upperStartCamelObject}}Model) Delete({{.lowerStartCamelPrimaryKey}} int64) error {
	result := m.db.Where("{{.originalPrimaryKey}} = ?", {{.lowerStartCamelPrimaryKey}}).Delete(&{{.upperStartCamelObject}}{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (m *default{{.upperStartCamelObject}}Model) Deletes({{.lowerStartCamelPrimaryKey}}s []int64) error {
	// 开启事务
	tx := m.db.Begin()
	// 循环遍历ids切片
	for _, id := range {{.lowerStartCamelPrimaryKey}}s {
		// 构建条件进行删除
		result := tx.Where("{{.originalPrimaryKey}} = ?", id).Delete(&{{.upperStartCamelObject}}{})
		if result.Error != nil {
			tx.Rollback() // 发生错误时回滚事务
			return result.Error
		}
	}
	// 提交事务
	err := tx.Commit().Error
	if err != nil {
		tx.Rollback() // 提交事务发生错误时回滚
		return err
	}
	return nil
}
