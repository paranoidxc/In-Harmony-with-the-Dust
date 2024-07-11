package autocurd

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"text/template"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/model"

	"github.com/go-cmd/cmd"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type CreateAutoCurdLogicGorm struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 新增
func NewCreateAutoCurdLogicGorm(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAutoCurdLogicGorm {
	return &CreateAutoCurdLogicGorm{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAutoCurdLogicGorm) CreateAutoCurd() error {
	// 获取结构体
	reqModelName := "TestGorm"
	var modelStruct interface{}
	for k, v := range model.AutoCrudModelList {
		if k == reqModelName {
			modelStruct = v
		}
	}
	m := reflect.TypeOf(modelStruct)
	// 大驼峰名字
	// 结构体名称
	name := m.Name()
	underlineNameGorm := GetUnderlineWord(name)
	// name = strings.Replace(name, "Tmp", "", -1)
	fmt.Println("结构体名称", name)

	underlineName := GetUnderlineWord(name)
	lowerCaseName := strings.ToLower(name[:1]) + name[1:]

	// 主键名字
	primaryKeyName := m.Field(0).Name
	primaryKeyJson := m.Field(0).Tag.Get("json")

	createStruct := ""
	createContent := ""
	deleteContentRequest := ""
	deleteContentResponse := ""
	deletesContentRequest := ""
	deletesContentResponse := ""
	updateContent := ""
	detailContentRequest := ""
	listContent := ""
	pageContent := ""
	// 前端字段kv
	vueFields := []map[string]string{}

	// struct需要的字段
	for i := 0; i < m.NumField(); i++ {
		field := m.Field(i)
		fmt.Printf("field %+v\n", field)
		item := fmt.Sprintf(`%v %v`, field.Name, field.Type)
		tag := `json:"` + field.Tag.Get("json") + `"`
		createStruct += (item + " `" + tag + "`" + "\n")
	}
	fmt.Println("createStruct", createStruct)

	// create需要的字段
	for i := 1; i < m.NumField(); i++ {
		field := m.Field(i)
		item := fmt.Sprintf(`%v %v`, field.Name, field.Type)
		// 过滤掉不需要的字段
		if filterName(field.Name, "create_request") {
			continue
		}
		tag := `json:"` + field.Tag.Get("json") + `"`

		validate := ""
		label := ""

		if len(field.Tag.Get("validate")) > 0 {
			validate = ` validate:"` + field.Tag.Get("validate") + `"`
		}

		if len(field.Tag.Get("label")) > 0 {
			label = ` label:"` + field.Tag.Get("label") + `"`
		}

		createContent += (item + " `" + tag + validate + label + "`" + "\n")
	}
	logx.Info("createContent", createContent)

	// delete Request需要的字段
	// 取第一个 所以第一个需要是主键
	for i := 0; i < 1; i++ {
		field := m.Field(i)
		item := fmt.Sprintf(`%v %v`, field.Name, field.Type)
		// 过滤掉不需要的字段
		if filterName(field.Name, "delete_request") {
			continue
		}
		tag := `json:"` + field.Tag.Get("json") + `"`
		deleteContentRequest += (item + " `" + tag + "`" + "\n")
	}
	// delete Response需要的字段
	for i := 0; i < 1; i++ {
		field := m.Field(i)
		item := fmt.Sprintf(`%v %v`, field.Name, field.Type)
		// 过滤掉不需要的字段
		if filterName(field.Name, "delete_response") {
			continue
		}
		tag := `json:"` + field.Tag.Get("json") + `"`
		deleteContentResponse += (item + " `" + tag + "`" + "\n")
	}

	// deletes Request需要的字段
	// 取第一个 所以第一个需要是主键
	for i := 0; i < 1; i++ {
		field := m.Field(i)
		item := fmt.Sprintf(`%v []%v`, field.Name, field.Type)
		// 过滤掉不需要的字段
		if filterName(field.Name, "deletes_request") {
			continue
		}
		tag := `json:"` + field.Tag.Get("json") + `"`
		deletesContentRequest += (item + " `" + tag + "`" + "\n")
	}
	// deletes Response需要的字段
	for i := 0; i < 1; i++ {
		field := m.Field(i)
		item := fmt.Sprintf(`%v %v`, field.Name, field.Type)
		// 过滤掉不需要的字段
		if filterName(field.Name, "deletes_response") {
			continue
		}
		tag := `json:"[]` + field.Tag.Get("json") + `"`
		deletesContentResponse += (item + " `" + tag + "`" + "\n")
	}
	// update及列表返回等需要的字段
	for i := 0; i < m.NumField(); i++ {
		field := m.Field(i)
		item := fmt.Sprintf(`%v %v`, field.Name, field.Type)
		// 过滤掉不需要的字段
		if filterName(field.Name, "update_request") {
			continue
		}
		tag := `json:"` + field.Tag.Get("json") + `"`
		updateContent += (item + " `" + tag + "`" + "\n")
	}

	// detail及列表返回等需要的字段
	for i := 0; i < 1; i++ {
		field := m.Field(i)
		item := fmt.Sprintf(`%v %v`, field.Name, field.Type)
		// 过滤掉不需要的字段
		if filterName(field.Name, "detail_request") {
			continue
		}
		tag := `form:"` + field.Tag.Get("json") + `"`
		detailContentRequest += (item + " `" + tag + "`" + "\n")
	}

	// list request需要的字段
	for i := 1; i < m.NumField(); i++ {
		field := m.Field(i)
		item := fmt.Sprintf(`%v %v`, field.Name, field.Type)
		// 过滤掉不需要的字段
		if filterName(field.Name, "list_request") {
			continue
		}
		tag := `form:"` + field.Tag.Get("json") + ",optional" + `"`
		listContent += (item + " `" + tag + "`" + "\n")
	}
	// page request需要的字段
	for i := 1; i < m.NumField(); i++ {
		field := m.Field(i)
		item := fmt.Sprintf(`%v %v`, field.Name, field.Type)
		// 过滤掉不需要的字段
		if filterName(field.Name, "page_request") {
			continue
		}
		tag := `form:"` + field.Tag.Get("json") + ",optional" + `"`
		pageContent += (item + " `" + tag + "`" + "\n")
	}

	for i := 1; i < m.NumField(); i++ {
		field := m.Field(i)
		// 过滤掉不需要的字段
		if filterName(field.Name, "vue_fields") {
			continue
		}
		key := field.Tag.Get("json")
		item := field.Tag.Get("gorm")
		label := getCommentFromGormTag(item)
		column := getColumnFromGormTag(item)
		fmt.Println("label", label)
		vueFields = append(vueFields, map[string]string{"Key": key, "Label": label, "Name": field.Name, "Column": column})
	}

	CreateStruct := getStruct(name, createStruct)
	CreateRequest := getCreateRequest(name, createContent)
	//CreateResponse := ""
	//CreateResponse := getCreateResponse(name, updateContent)
	DeleteRequest := getDeleteRequest(name, deleteContentRequest)
	//DeleteResponse := ""
	//DeleteResponse := getDeleteResponse(name, deleteContentResponse)
	DeletesRequest := getDeletesRequest(name, deletesContentRequest)
	//DeletesResponse := ""
	//DeletesResponse := getDeletesResponse(name, deletesContentResponse)
	UpdateRequest := getUpdateRequest(name, updateContent)
	//UpdateResponse := ""
	//UpdateResponse := getUpdateResponse(name, updateContent)
	DetailRequest := getDetailRequest(name, detailContentRequest)
	DetailResponse := getDetailResponse(name, updateContent)
	ListRequest := getListRequest(name, listContent)
	ListResponse := getListResponse(name)
	PageRequest := getPageRequest(name, pageContent)
	PageResponse := getPageResponse(name)
	ServerContent := getServerContent(name, underlineName, lowerCaseName)

	res := ""
	res += CreateStruct + "\n"
	res += CreateRequest + "\n"
	//res += CreateResponse + "\n"
	res += DeleteRequest + "\n"
	//res += DeleteResponse + "\n"
	res += DeletesRequest + "\n"
	//res += DeletesResponse + "\n"
	res += UpdateRequest + "\n"
	//res += UpdateResponse + "\n"
	res += DetailRequest + "\n"
	res += DetailResponse + "\n"
	res += ListRequest + "\n"
	res += ListResponse + "\n"
	res += PageRequest + "\n"
	res += PageResponse + "\n"
	res += ServerContent + "\n"

	//时间字段替换为 string
	res = strings.Replace(res, "time.Time", "string", -1)
	res = strings.Replace(res, "gorm.DeletedAt", "string", -1)

	// 生成.api文件
	fmt.Println(primaryKeyName)
	err := createApiFile(underlineName, res)
	// 将.api文件加入总的api文件
	err = addApiFile(underlineName)
	// 运行goctl命令生成代码
	err = goCtlGenFileGorm()
	err = goCtlGenModelFileGorm(l, underlineNameGorm, modelStruct)
	//编辑模型文件
	err = editModelFileGorm(name, underlineName)
	// 编辑logic文件
	err = editLogicFileGorm(name, underlineName, primaryKeyName, primaryKeyJson, vueFields)
	// 生成前端文件
	err = genWebApiFile(underlineName, lowerCaseName, primaryKeyJson)
	err = genWebVueFileGorm(name, underlineName, primaryKeyJson, vueFields)
	//resp = &types.Empty{}
	if err != nil {
		fmt.Println("出错了：", err)

	}

	return err
}

// 判断名称是否在过滤的字段中
func filterName(name string, typename string) bool {
	filterArray := map[string][]string{
		"create":  []string{"CreatedAt", "UpdatedAt", "DeletedAt"},
		"default": []string{"CreatedAt", "UpdatedAt", "DeletedAt"},
	}
	//如果filterArray[typename] 不存在 取 default配置
	var list []string
	list = filterArray[typename]
	if list == nil {
		list = filterArray["default"]
	}
	for _, v := range list {
		if v == name {
			return true
		}
	}
	return false
}
func genWebVueFileGorm(name, underlineName, primaryKeyJson string, vueFields interface{}) error {
	projectWd, _ := os.Getwd()
	fileDir := filepath.Join(projectWd, "../../../../web")
	filePath := filepath.Join(projectWd, "../tpl/table.tpl")
	tpl, err := template.ParseFiles(filePath)
	if err != nil {
		fmt.Println("create template failed, err:", err)
		return err
	}
	file, err := os.Create(fileDir + "/src/views/feat/" + underlineName + ".vue")
	defer file.Close()
	data := map[string]interface{}{
		"Name":           name,
		"UnderlineName":  underlineName,
		"PrimaryKeyJson": primaryKeyJson,
		"VueFields":      vueFields,
	}
	err = tpl.Execute(file, data)
	if err != nil {
		return err
	}
	return nil
}
func goCtlGenFileGorm() error {
	// goctl生成文件
	cmdArgs := []string{"api", "go", "-api", "core.api", "--style", "goZero", "-dir", ".", "--home", "../tpl_gorm"}
	c := cmd.NewCmd("goctl", cmdArgs...)
	<-c.Start()
	fmt.Println("goctl生成逻辑文件")
	return nil
}

// 编辑模型
func editModelFileGorm(name, underlineName string) error {
	// 读取文件内容
	projectWd, _ := os.Getwd()
	tmpfile := filepath.Join(projectWd, "../../model/tmpAutoFeat/"+underlineName+".go")
	tmpcontent, err := ioutil.ReadFile(tmpfile)
	if err != nil {
		fmt.Println("读取初始模型文件失败:", err)
	}

	// 定义正则表达式来匹配结构体字段定义部分
	re := regexp.MustCompile(`type\s+(\w+)\s+struct\s+{([\s\S]*?)\}`)
	var fields string
	// 寻找匹配的结构体定义并提取字段部分
	matches := re.FindAllStringSubmatch(string(tmpcontent), -1)
	for _, match := range matches {
		structName := match[1]
		fields = match[2]
		fmt.Println("Structure Name:", structName)
		fmt.Println(fields)
	}
	// 读取已生成的模型内容
	fileName := strings.ToLower(name)
	modelFile := filepath.Join(projectWd, "../../model/", fileName+"Model_gen.go")
	content, err := ioutil.ReadFile(modelFile) // 读取文件内容
	if err != nil {
		fmt.Println("读取生成模型文件失败:", err)
	}
	modifiedContent := string(content)
	modifiedContent = strings.Replace(modifiedContent, "replace string", fields, -1)
	// 将修改后的内容写回文件
	err = ioutil.WriteFile(modelFile, []byte(modifiedContent), 0644)
	if err != nil {
		fmt.Println("无法写入文件:", err)
		return err
	}
	fmt.Println("编辑模型文件")
	return nil
}
func goCtlGenModelFileGorm(l *CreateAutoCurdLogicGorm, underlineName string, model interface{}) error {

	// goctl生成文件
	//url := `-url=" + l.svcCtx.Config.Mysql.DataSource
	url := strings.Replace(l.svcCtx.Config.Mysql.DataSource, "?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai", "", -1)

	//gorm连接
	db, err := gorm.Open(mysql.Open(l.svcCtx.Config.Mysql.DataSource), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   false,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
	})
	if err != nil {
		panic("连接数据库失败")
	}
	//使用AutoMigrate自动生成表
	if err := db.AutoMigrate(
		&model,
	); err != nil {
		panic(err.Error())
	}
	//table := "-table=\"" + underlineName + "\""
	//cmdArgs := []string{"model", "mysql", "datasource", url, table, `-dir="../../model"`, "--style", "goZero", "--home", "../tpl"}
	//cmdArgs := []string{"model", "mysql", "datasource", url, table, "-dir", ".", "-cache", "true", "--style", "goZero", "--home", "../tpl"}
	cmdArgs := []string{"model", "mysql", "datasource", "-url", url, "-table", underlineName, "-dir", "../../model", "-cache", "true", "--style", "goZero", "--home", "../tpl_gorm"}
	fmt.Println("go model:", strings.Join(cmdArgs, " "))
	c := cmd.NewCmd("goctl", cmdArgs...)
	fmt.Println("goctl生成模型文件")
	<-c.Start()
	return nil
}
func editLogicFileGorm(name, underlineName, primaryKeyName, primaryKeyJson string, vueFields []map[string]string) error {
	// 新增逻辑
	createLogic := fmt.Sprintf(`
func (l *%vCreateLogic) %vCreate(req *types.%vCreateReq) (err error) {
	var modelParams = new(model.%v)
	err = copier.Copy(modelParams, req)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}
	err = l.svcCtx.Feat%vModel.Insert(modelParams)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	return
}
`, name, name, name, name, name)

	// 删除逻辑
	deleteLogic := fmt.Sprintf(`
func (l *%vDeleteLogic) %vDelete(req *types.%vDeleteReq) (err error) {
	err = l.svcCtx.Feat%vModel.Delete(req.%v)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	return
}
`, name, name, name, name, primaryKeyName)

	deletesLogic := fmt.Sprintf(`
func (l *%vDeletesLogic) %vDeletes(req *types.%vDeletesReq) (err error) {
	if len(req.%v) > 0  {
		err = l.svcCtx.Feat%vModel.Deletes(req.%v)
		if err != nil {
			return  errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}
	} else {
		return errorx2.NewSystemError(errorx2.ParamErrorCode, err.Error())
	}

	return
}
`, name, name, name, primaryKeyName, name, primaryKeyName)

	// 修改逻辑
	updateLogic := fmt.Sprintf(`
func (l *%vUpdateLogic) %vUpdate(req *types.%vUpdateReq) (err error) {
	modelParams := &model.%v{}
	modelParams, err = l.svcCtx.Feat%vModel.FindOne(req.%v)
	if err != nil {
		return errorx2.NewDefaultError(errorx2.UserIdErrorCode)
	}

	err = copier.Copy(modelParams, req)
	if err != nil {
		logx.Error("复制参数失败", err)
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	err = l.svcCtx.Feat%vModel.Update(modelParams)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	return
}
`, name, name, name, name, name, primaryKeyName, name)

	// 详情逻辑
	detailLogic := fmt.Sprintf(`
func (l *%vDetailLogic) %vDetail(req *types.%vDetailReq) (resp *types.%vDetailResp, err error) {
	resp = &types.%vDetailResp{}
	item := &model.%v{}
	item, err = l.svcCtx.Feat%vModel.FindOne(req.%v)
	err = copier.Copy(resp, item)
	if err != nil {
		logx.Error("复制结果失败", err)
		return nil, err
	}
	return
}
`, name, name, name, name, name, name, name, primaryKeyName)

	// 列表逻辑
	listLogic, _ := getListLogicGorm(name, vueFields)

	// 分页列表逻辑
	pageLogic, _ := getPageLogicGorm(name, vueFields)

	// 生成存放的文件路径
	fileName := strings.ToLower(name)
	projectWd, _ := os.Getwd()
	createLogicFile := filepath.Join(projectWd, "./internal/logic/feat/"+underlineName+"/", fileName+"CreateLogic.go")
	deleteLogicFile := filepath.Join(projectWd, "./internal/logic/feat/"+underlineName+"/", fileName+"DeleteLogic.go")
	deletesLogicFile := filepath.Join(projectWd, "./internal/logic/feat/"+underlineName+"/", fileName+"DeletesLogic.go")
	updateLogicFile := filepath.Join(projectWd, "./internal/logic/feat/"+underlineName+"/", fileName+"UpdateLogic.go")
	detailLogicFile := filepath.Join(projectWd, "./internal/logic/feat/"+underlineName+"/", fileName+"DetailLogic.go")
	listLogicFile := filepath.Join(projectWd, "./internal/logic/feat/"+underlineName+"/", fileName+"ListLogic.go")
	pageLogicFile := filepath.Join(projectWd, "./internal/logic/feat/"+underlineName+"/", fileName+"PageLogic.go")

	fileList := map[string]string{
		createLogicFile:  createLogic,
		deleteLogicFile:  deleteLogic,
		deletesLogicFile: deletesLogic,
		updateLogicFile:  updateLogic,
		detailLogicFile:  detailLogic,
		listLogicFile:    listLogic,
		pageLogicFile:    pageLogic,
	}
	for k, v := range fileList {
		content, err := ioutil.ReadFile(k) // 读取文件内容
		if err != nil {
			fmt.Printf("读取文件%s失败：%v\n", k, err)
			return err
		}
		tmpMethods := []string{
			"Create", "Update", "Detail", "Deletes", "Delete", "List", "Page",
		}
		method := ""
		for _, _method := range tmpMethods {
			contains := strings.Contains(k, _method)
			if contains {
				method = _method
				break
			}
		}
		methods := []string{
			method,
		}

		modifiedContent := string(content)
		for _, method := range methods {
			pattern := fmt.Sprintf(`func \(l \*%v%vLogic.* \{[\s\S]*?\}`, name, method)
			regex := regexp.MustCompile(pattern)
			modifiedContent = regex.ReplaceAllString(modifiedContent, "")

			// 正则表达式模式
			//pattern = `(\r?\n){4}`
			//modifiedContent = regexp.MustCompile(pattern).ReplaceAllString(modifiedContent, "")
			pattern = `(\r?\n){3}`
			modifiedContent = regexp.MustCompile(pattern).ReplaceAllString(modifiedContent, "")
		}
		if method == "Delete" {
			modifiedContent = strings.Replace(modifiedContent, `"github.com/jinzhu/copier"`, "", -1)
			modifiedContent = strings.Replace(modifiedContent, `"zero-zone/app/model"`, "", -1)
		}
		if method == "Deletes" {
			modifiedContent = strings.Replace(modifiedContent, `"github.com/jinzhu/copier"`, "", -1)
			modifiedContent = strings.Replace(modifiedContent, `"zero-zone/app/model"`, "", -1)
		}
		if method == "Detail" {
			modifiedContent = strings.Replace(modifiedContent, `errorx2 "zero-zone/pkg/errorx"`, "", -1)
		}
		if method == "List" {
			modifiedContent = strings.Replace(modifiedContent, `"zero-zone/app/model"`, "", -1)
		}
		if method == "Page" {
			modifiedContent = strings.Replace(modifiedContent, `"zero-zone/app/model"`, "", -1)
		}

		// 将修改后的内容写回文件
		err = ioutil.WriteFile(k, []byte(modifiedContent), 0644)
		if err != nil {
			fmt.Println("无法写入文件:", err)
			return err
		}

		// 向文件中追加内容
		itemFileRaw, err := os.OpenFile(k, os.O_APPEND|os.O_WRONLY, 0644) // 打开文件并设置写入模式
		if err != nil {
			fmt.Println("打开index.api文件失败：", err)
			return err
		}
		defer itemFileRaw.Close() // 关闭文件
		_, err = itemFileRaw.WriteString(fmt.Sprintf("\n%v\n", v))
		if err != nil {
			fmt.Println("向index.api文件中追加内容失败：", err)
			return err
		}
	}
	fmt.Println("编辑逻辑文件")

	return nil
}
func getListLogicGorm(name string, vueFields []map[string]string) (string, error) {
	projectWd, _ := os.Getwd()
	file := filepath.Join(projectWd, "../tpl_gorm/list.tpl")
	tpl, err := template.ParseFiles(file)
	if err != nil {
		fmt.Println("create template failed, err:", err)
		return "", err
	}
	var tplBuffer bytes.Buffer
	data := map[string]interface{}{
		"Name":      name,
		"VueFields": vueFields,
	}
	err = tpl.Execute(&tplBuffer, data)
	if err != nil {
		return "", err
	}
	return tplBuffer.String(), nil
}

func getPageLogicGorm(name string, vueFields []map[string]string) (string, error) {
	projectWd, _ := os.Getwd()
	file := filepath.Join(projectWd, "../tpl_gorm/page.tpl")
	tpl, err := template.ParseFiles(file)
	if err != nil {
		fmt.Println("create template failed, err:", err)
		return "", err
	}
	var tplBuffer bytes.Buffer
	data := map[string]interface{}{
		"Name":      name,
		"VueFields": vueFields,
	}
	err = tpl.Execute(&tplBuffer, data)
	if err != nil {
		return "", err
	}
	return tplBuffer.String(), nil
}
