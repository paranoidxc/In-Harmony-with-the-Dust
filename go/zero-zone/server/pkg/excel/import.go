package excel

import (
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"reflect"
	"strconv"
	"strings"
	"zero-zone/applet/pkg/excel/model"
)

// ImportExcel 导入数据（单个sheet）
// 需要在传入的结构体中的字段加上tag：excel:"title:列头名称;"
// f 获取到的excel对象、dst 导入目标对象【传指针】
// headIndex 表头（最后一行表头）的索引，从0开始（用于获取表头名字）
// startRow 标题+表头行行数（从第startRow+1行开始扫）
func ImportExcel(f *excelize.File, dst interface{}, headIndex, startRow int) (err error) {
	sheetName := f.GetSheetName(0) // 单个sheet时，默认读取第一个sheet
	err = importData(f, dst, sheetName, headIndex, startRow)
	return
}

// ImportBySheet 导入数据（读取指定sheet）sheetName Sheet名称
func ImportBySheet(f *excelize.File, dst interface{}, sheetName string, headIndex, startRow int) (err error) {
	// 当需要读取多个sheet时，可以通过下面的方式，来调用 ImportBySheet 这个函数
	/*sheetList := f.GetSheetList()
	for _, sheetName := range sheetList {
		ImportBySheet(f,dst,sheetName,headIndex,startRow)
	}*/
	err = importData(f, dst, sheetName, headIndex, startRow)
	return
}

// 解析数据
func importData(f *excelize.File, dst interface{}, sheetName string, headIndex, startRow int) (err error) {
	rows, err := f.GetRows(sheetName) // 获取所有行
	if err != nil {
		err = errors.New(sheetName + "工作表不存在")
		return
	}
	dataValue := reflect.ValueOf(dst) // 取目标对象的元素类型、字段类型和 tag
	// 判断数据的类型
	if dataValue.Kind() != reflect.Ptr || dataValue.Elem().Kind() != reflect.Slice {
		err = errors.New("无效的数据类型")
	}
	heads := []string{}                        // 表头
	dataType := dataValue.Elem().Type().Elem() // 获取导入目标对象的类型信息
	// 遍历行，解析数据并填充到目标对象中
	for rowIndex, row := range rows {
		if rowIndex == headIndex {
			heads = row
		}
		if rowIndex < startRow { // 跳过头行
			continue
		}
		newData := reflect.New(dataType).Elem() // 创建新的目标对象
		// 遍历目标对象的字段
		for i := 0; i < dataType.NumField(); i++ {
			// 这里要用构造函数，构造函数里指定了Index默认值为-1，当目标结构体的tag没有指定index的话，那么 excelTag.Index 就一直为0
			// 那么 row[excelizeIndex] 就始终是 row[0]，始终拿的是第一列的数据
			var excelTag = model.NewExcelTag()
			field := dataType.Field(i) // 获取字段信息和tag
			tag := field.Tag.Get(model.ExcelTagKey)
			if tag == "" { // 如果tag不存在，则跳过
				continue
			}
			err = excelTag.GetTag(tag)
			if err != nil {
				return
			}
			cellValue := ""
			if excelTag.Index >= 0 { // 当tag里指定了index时，根据这个index来拿数据
				excelizeIndex := excelTag.Index // 解析tag的值
				if excelizeIndex >= len(row) {  // 防止下标越界
					continue
				}
				cellValue = row[excelizeIndex] // 获取单元格的值
			} else { // 否则根据表头名称来拿数据
				if model.IsContain(heads, excelTag.Name) { // 当tag里的表头名称和excel表格里面的表头名称相匹配时
					if i >= len(row) { // 防止下标越界
						continue
					}
					cellValue = row[i] // 获取单元格的值
				}
			}
			// 替换
			if excelTag.Replace != "" {
				split := strings.Split(excelTag.Replace, ",")
				for j := range split {
					s := strings.Split(split[j], "_") // 根据下划线进行分割，格式：需要替换的内容_替换后的内容
					// 导出时是取下划线后面的值，导入则是反过来取下划线前面的值
					if s[1] == cellValue {
						cellValue = s[0]
					}
				}
			}
			// 根据字段类型设置值
			switch field.Type.Kind() {
			case reflect.Int:
				v, _ := strconv.ParseInt(cellValue, 10, 64)
				newData.Field(i).SetInt(v)
			case reflect.Bool:
				v, _ := strconv.ParseBool(cellValue)
				newData.Field(i).SetBool(v)
			case reflect.String:
				newData.Field(i).SetString(cellValue)
			}
		}
		// 将新的目标对象添加到导入目标对象的slice中
		dataValue.Elem().Set(reflect.Append(dataValue.Elem(), newData))
	}
	return
}
