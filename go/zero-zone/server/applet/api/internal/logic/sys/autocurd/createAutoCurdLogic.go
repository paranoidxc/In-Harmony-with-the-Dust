package autocurd

import (
	"bytes"
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"text/template"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
	"zero-zone/applet/model"

	"github.com/go-cmd/cmd"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAutoCurdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 新增
func NewCreateAutoCurdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAutoCurdLogic {
	return &CreateAutoCurdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAutoCurdLogic) CreateAutoCurd(req *types.AutoCurdCreateReq) error {
	// 获取结构体
	reqModelName := req.Name
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
	name = strings.Replace(name, "Tmp", "", -1)
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
	detailContentResponse := ""
	listContent := ""
	pageContent := ""
	// 前端字段kv
	vueFields := []map[string]string{}
	minVueFields := []map[string]string{}

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
		tag := `json:"` + field.Tag.Get("json") + `,optional"`

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
		tag := `json:"` + field.Tag.Get("json") + `"`
		deleteContentRequest += (item + " `" + tag + "`" + "\n")
	}
	// delete Response需要的字段
	for i := 0; i < 1; i++ {
		field := m.Field(i)
		item := fmt.Sprintf(`%v %v`, field.Name, field.Type)
		tag := `json:"` + field.Tag.Get("json") + `"`
		deleteContentResponse += (item + " `" + tag + "`" + "\n")
	}

	// deletes Request需要的字段
	// 取第一个 所以第一个需要是主键
	for i := 0; i < 1; i++ {
		field := m.Field(i)
		item := fmt.Sprintf(`%v []%v`, field.Name, field.Type)
		tag := `json:"` + field.Tag.Get("json") + `"`
		deletesContentRequest += (item + " `" + tag + "`" + "\n")
	}
	// deletes Response需要的字段
	for i := 0; i < 1; i++ {
		field := m.Field(i)
		item := fmt.Sprintf(`%v %v`, field.Name, field.Type)
		tag := `json:"[]` + field.Tag.Get("json") + `"`
		deletesContentResponse += (item + " `" + tag + "`" + "\n")
	}
	// update及列表返回等需要的字段
	for i := 0; i < m.NumField(); i++ {
		field := m.Field(i)
		item := fmt.Sprintf(`%v %v`, field.Name, field.Type)
		tag := ""
		if i == 0 {
			tag = `json:"` + field.Tag.Get("json") + `"`
		} else {
			tag = `json:"` + field.Tag.Get("json") + `,optional"`
		}
		validate := ""
		label := ""

		if len(field.Tag.Get("validate")) > 0 {
			validate = ` validate:"` + field.Tag.Get("validate") + `"`
		}

		if len(field.Tag.Get("label")) > 0 {
			label = ` label:"` + field.Tag.Get("label") + `"`
		}
		updateContent += (item + " `" + tag + validate + label + "`" + "\n")
	}

	for i := 0; i < m.NumField(); i++ {
		field := m.Field(i)
		item := fmt.Sprintf(`%v %v`, field.Name, field.Type)
		tag := `json:"` + field.Tag.Get("json") + `"`
		label := ""
		if len(field.Tag.Get("label")) > 0 {
			label = ` label:"` + field.Tag.Get("label") + `"`
		}
		detailContentResponse += (item + " `" + tag + label + "`" + "\n")
	}

	// detail及列表返回等需要的字段
	for i := 0; i < 1; i++ {
		field := m.Field(i)
		item := fmt.Sprintf(`%v %v`, field.Name, field.Type)
		tag := `form:"` + field.Tag.Get("json") + `"`
		detailContentRequest += (item + " `" + tag + "`" + "\n")
	}

	// list request需要的字段
	for i := 1; i < m.NumField(); i++ {
		field := m.Field(i)
		item := fmt.Sprintf(`%v %v`, field.Name, field.Type)
		tag := `form:"` + field.Tag.Get("json") + ",optional" + `"`
		listContent += (item + " `" + tag + "`" + "\n")
	}
	listContent += fmt.Sprintf(`%v %v`, "IncludeDeleted", "int")
	listtag := `form:"includeDeleted,optional"`
	listContent += " `" + listtag + "`" + "\n"

	// page request需要的字段
	for i := 1; i < m.NumField(); i++ {
		field := m.Field(i)
		item := fmt.Sprintf(`%v %v`, field.Name, field.Type)
		tag := `form:"` + field.Tag.Get("json") + ",optional" + `"`
		pageContent += (item + " `" + tag + "`" + "\n")
	}
	pageContent += fmt.Sprintf(`%v %v`, "IncludeDeleted", "int")
	pagetag := `form:"includeDeleted,optional"`
	pageContent += " `" + pagetag + "`" + "\n"

	for i := 1; i < m.NumField(); i++ {
		field := m.Field(i)
		key := field.Tag.Get("json")
		item := field.Tag.Get("db")
		label := field.Tag.Get("label")
		comment := getCommentFromGormTag(item)
		if comment == "label" || comment == "" {
			comment = label
		}
		column := getColumnFromGormTag(item)
		fmt.Println("comment", comment)
		fmt.Println("key", key)
		vueFields = append(vueFields, map[string]string{"Key": key, "Label": comment, "Name": field.Name, "Column": column})

		if !(strings.ToLower(key) == "createdat" || strings.ToLower(key) == "updatedat" || strings.ToLower(key) == "deletedat") {
			minVueFields = append(minVueFields, map[string]string{"Key": key, "Label": comment, "Name": field.Name, "Column": column})
		}
	}

	CreateStruct := getStruct(name, createStruct)
	CreateRequest := getCreateRequest(name, createContent)
	DeleteRequest := getDeleteRequest(name, deleteContentRequest)
	DeletesRequest := getDeletesRequest(name, deletesContentRequest)
	UpdateRequest := getUpdateRequest(name, updateContent)
	DetailRequest := getDetailRequest(name, detailContentRequest)
	DetailResponse := getDetailResponse(name, detailContentResponse)
	ListRequest := getListRequest(name, listContent)
	ListResponse := getListResponse(name)
	PageRequest := getPageRequest(name, pageContent)
	PageResponse := getPageResponse(name)
	ServerContent := getServerContent(name, underlineName, lowerCaseName)

	res := ""
	res += CreateStruct + "\n"
	res += CreateRequest + "\n"
	res += DeleteRequest + "\n"
	res += DeletesRequest + "\n"
	res += UpdateRequest + "\n"
	res += DetailRequest + "\n"
	res += DetailResponse + "\n"
	res += ListRequest + "\n"
	res += ListResponse + "\n"
	res += PageRequest + "\n"
	res += PageResponse + "\n"
	res += ServerContent + "\n"

	isApiFile := false
	isHandleFile := false
	isLogicWrite := false
	isModelFile := false
	isVueFile := false
	isMenu := false
	if req.IsAll == 1 {
		isApiFile = true
		isHandleFile = true
		isLogicWrite = true
		isModelFile = true
		isVueFile = true
		isMenu = true
	}

	if req.IsApi == 1 {
		isApiFile = true
	}

	if req.IsHandle == 1 {
		isHandleFile = true
	}

	if req.IsLogicWrite == 1 {
		isLogicWrite = true
	}

	if req.IsModel == 1 {
		isModelFile = true
	}

	if req.IsVue == 1 {
		isVueFile = true
	}

	if req.IsMenu == 1 {
		isMenu = true
	}

	// 生成.api文件
	fmt.Println(primaryKeyName)
	var err error
	if isApiFile {
		err = createApiFile(underlineName, res)
		// 将.api文件加入总的api文件
		err = addApiFile(underlineName)
	}
	// 运行goctl命令生成代码
	if isHandleFile {
		err = goCtlGenFile()
		if err != nil {
			fmt.Println("出错了：", err)
		}
	}

	if isModelFile {
		err = goCtlGenModelFile(l, underlineName)
		if err != nil {
			fmt.Println("出错了：", err)
		}
	}

	if isLogicWrite {
		// 编辑logic文件
		err = editLogicFile(name, underlineName, primaryKeyName, primaryKeyJson, vueFields)
		if err != nil {
			fmt.Println("出错了：", err)
		}
	}

	// 生成前端文件
	if isVueFile {
		err = genWebApiFile(underlineName, lowerCaseName, primaryKeyJson)
		if err != nil {
			fmt.Println("出错了：", err)
		}
		err = genWebVueFile(name, underlineName, primaryKeyJson, vueFields, minVueFields)
		if err != nil {
			fmt.Println("出错了：", err)
		}
	}

	if isMenu {
		_, err := l.svcCtx.SysPermMenuModel.FindOneByName(l.ctx, reqModelName)
		if err != nil && err == sqlx.ErrNotFound {
			res, err := l.svcCtx.SysPermMenuModel.Insert(l.ctx, &model.SysPermMenu{
				ParentId: 42,
				Name:     reqModelName,
				Router:   "/feat/" + underlineName,
				Perms:    "[]",
				Type:     1,
				ViewPath: "views/feat/" + underlineName,
				IsShow:   1,
			})
			if err != nil {
				fmt.Println("出错了1：", err)
			}
			menuId, err := res.LastInsertId()
			if err != nil {
				fmt.Println("出错了2：", err)
			}
			//menu, _ := l.svcCtx.SysPermMenuModel.FindOneByName(l.ctx, reqModelName)
			l.svcCtx.SysPermMenuModel.Insert(l.ctx, &model.SysPermMenu{
				//ParentId: menu.Id,
				ParentId: menuId,
				Name:     reqModelName + "查询列表",
				Router:   "",
				Perms:    "[\"admin/feat/" + reqModelName + "/page\"]",
				Type:     2,
				ViewPath: "",
				IsShow:   1,
			})
			l.svcCtx.SysPermMenuModel.Insert(l.ctx, &model.SysPermMenu{
				//ParentId: menu.Id,
				ParentId: menuId,
				Name:     reqModelName + "所有列表",
				Router:   "",
				Perms:    "[\"admin/feat/" + reqModelName + "/list\"]",
				Type:     2,
				ViewPath: "",
				IsShow:   1,
			})
			l.svcCtx.SysPermMenuModel.Insert(l.ctx, &model.SysPermMenu{
				//ParentId: menu.Id,
				ParentId: menuId,
				Name:     reqModelName + "新建",
				Router:   "",
				Perms:    "[\"admin/feat/" + reqModelName + "/create\"]",
				Type:     2,
				ViewPath: "",
				IsShow:   1,
			})
			l.svcCtx.SysPermMenuModel.Insert(l.ctx, &model.SysPermMenu{
				//ParentId: menu.Id,
				ParentId: menuId,
				Name:     reqModelName + "更新",
				Router:   "",
				Perms:    "[\"admin/feat/" + reqModelName + "/update\"]",
				Type:     2,
				ViewPath: "",
				IsShow:   1,
			})
			l.svcCtx.SysPermMenuModel.Insert(l.ctx, &model.SysPermMenu{
				//ParentId: menu.Id,
				ParentId: menuId,
				Name:     reqModelName + "详情",
				Router:   "",
				Perms:    "[\"admin/feat/" + reqModelName + "/detail\"]",
				Type:     2,
				ViewPath: "",
				IsShow:   1,
			})
			l.svcCtx.SysPermMenuModel.Insert(l.ctx, &model.SysPermMenu{
				//ParentId: menu.Id,
				ParentId: menuId,
				Name:     reqModelName + "删除",
				Router:   "",
				Perms:    "[\"admin/feat/" + reqModelName + "/delete\"]",
				Type:     2,
				ViewPath: "",
				IsShow:   1,
			})
			l.svcCtx.SysPermMenuModel.Insert(l.ctx, &model.SysPermMenu{
				//ParentId: menu.Id,
				ParentId: menuId,
				Name:     reqModelName + "批量删除",
				Router:   "",
				Perms:    "[\"admin/feat/" + reqModelName + "/deletes\"]",
				Type:     2,
				ViewPath: "",
				IsShow:   1,
			})
		}
	}

	return err
}

func getStruct(name, content string) string {
	return fmt.Sprintf(`type %v {
		%v
	}`, name, content)
}

func getCreateRequest(name, content string) string {
	return fmt.Sprintf(`type %vCreateReq {
		%v
	}`, name, content)
}
func getCreateResponse(name, content string) string {
	return fmt.Sprintf(`type %vCreateResp {
		%v
	}`, name, content)
}
func getDeleteRequest(name, content string) string {
	return fmt.Sprintf(`type %vDeleteReq {
		%v
	}`, name, content)
}
func getDeleteResponse(name, content string) string {
	return fmt.Sprintf(`type %vDeleteResp {
		%v
	}`, name, content)
}
func getDeletesRequest(name, content string) string {
	return fmt.Sprintf(`type %vDeletesReq {
		%v
	}`, name, content)
}
func getDeletesResponse(name, content string) string {
	return fmt.Sprintf(`type %vDeletesResp {
		%v
	}`, name, content)
}
func getUpdateRequest(name, content string) string {
	return fmt.Sprintf(`type %vUpdateReq {
		%v
	}`, name, content)
}
func getUpdateResponse(name, content string) string {
	return fmt.Sprintf(`type %vUpdateResp {
		%v
	}`, name, content)
}

func getDetailRequest(name, content string) string {
	return fmt.Sprintf(`type %vDetailReq {
		%v
	}`, name, content)
}
func getDetailResponse(name, content string) string {
	return fmt.Sprintf(`type %vDetailResp {
		%v
	}`, name, content)
}

// 列表
func getListRequest(name, listContent string) string {
	return fmt.Sprintf(`type %vListReq {
		%v
	}`, name, listContent)
}
func getListResponse(name string) string {
	tag_list := "`" + "json:" + `"list"` + "`"
	tag_total := "`" + "json:" + `"total"` + "`"
	return fmt.Sprintf(`type %vListResp {
		List  []%v %v
		Total int64   %v                         
	}`, name, name, tag_list, tag_total)
}

// 分页列表
func getPageRequest(name, pageContent string) string {
	return fmt.Sprintf(`type %vPageReq {
		PageReq
		%v
	}`, name, pageContent)
}
func getPageResponse(name string) string {
	tag_list := "`" + "json:" + `"list"` + "`"
	tag_pagination := "`" + "json:" + `"pagination"` + "`"
	return fmt.Sprintf(`type %vPageResp {
		List  []%v %v
		Pagination Pagination   %v                         
	}`, name, name, tag_list, tag_pagination)
}

func getServerContent(name, underlineName, lowerCaseName string) string {

	return fmt.Sprintf(`

@server(
	group: feat/%v
	prefix: /admin/feat/%v
	jwt:        JwtAuth
)

service core-api {
	@handler %vList
	get /list(%vListReq) returns (%vListResp)

	@handler %vPage
	get /page(%vPageReq) returns (%vPageResp)

	@handler %vCreate
	post /create(%vCreateReq)
	
	@handler %vDelete
	post /delete(%vDeleteReq)
	
	@handler %vDeletes
	post /deletes(%vDeletesReq)

	@handler %vUpdate
	post /update(%vUpdateReq)
	
	@handler %vDetail
	get /detail(%vDetailReq) returns (%vDetailResp)
}

`, underlineName, lowerCaseName, name, name, name, name, name, name, name, name, name, name, name, name, name, name, name, name, name)
}

func createApiFile(underlineName, res string) error {
	// 生成存放的文件路径
	fileName := underlineName + ".api"
	projectWd, _ := os.Getwd()
	fmt.Println("projectWd", projectWd)

	filePath := filepath.Join(projectWd, "/desc", fileName)
	fmt.Println(filePath)
	// 如果存在，就删除文件
	_, err := os.Stat(filePath)
	if err == nil {
		err := os.Remove(filePath)
		if err != nil {
			return err
		}
	}

	// 写入文件
	err = ioutil.WriteFile(filePath, []byte(res), 0644)
	// 格式化文件
	cmdArgs := []string{"api", "format", "-dir", filePath}
	c := cmd.NewCmd("goctl", cmdArgs...)
	<-c.Start()
	return nil
}

func addApiFile(underlineName string) error {
	coreApiFileName := "core.api"
	projectWd, _ := os.Getwd()
	filePath := filepath.Join(projectWd, "/"+coreApiFileName)

	fileContent, err := ioutil.ReadFile(filePath) // 读取文件内容
	if err != nil {

	}
	subStr := fmt.Sprintf(`"desc/%v.api"`, underlineName)
	if !strings.Contains(string(fileContent), subStr) {
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644) // 打开文件并设置写入模式
		if err != nil {
			fmt.Println("打开"+coreApiFileName+"文件失败：", err)
			return err
		}
		defer file.Close() // 关闭文件

		// 向文件中追加内容
		content := fmt.Sprintf(`%vimport "desc/%v.api"`, "\n", underlineName)
		_, err = file.WriteString(content)
		if err != nil {
			fmt.Println("向"+coreApiFileName+"文件中追加内容失败：", err)
			return err
		}
	}
	return nil
}

func goCtlGenFile() error {
	// goctl生成文件
	cmdArgs := []string{"api", "go", "-api", "core.api", "--style", "goZero", "-dir", ".", "--home", "../tpl"}
	c := cmd.NewCmd("goctl", cmdArgs...)
	<-c.Start()
	fmt.Println("goctl生成逻辑文件")
	return nil
}

func goCtlGenModelFile(l *CreateAutoCurdLogic, underlineName string) error {
	// goctl生成文件
	//url := `-url=" + l.svcCtx.Config.Mysql.DataSource
	url := strings.Replace(l.svcCtx.Config.Mysql.DataSource, "?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai", "", -1)
	//table := "-table=\"" + underlineName + "\""
	//cmdArgs := []string{"model", "mysql", "datasource", url, table, `-dir="../../model"`, "--style", "goZero", "--home", "../tpl"}
	//cmdArgs := []string{"model", "mysql", "datasource", url, table, "-dir", ".", "-cache", "true", "--style", "goZero", "--home", "../tpl"}
	cmdArgs := []string{"model", "mysql", "datasource", "-url", url, "-table", underlineName, "-dir", "../model", "-cache", "true", "--style", "goZero", "--home", "../tpl", "-i", "created_at,updated_at,deleted_at"}
	fmt.Println("go model:", "goctl "+strings.Join(cmdArgs, " "))
	c := cmd.NewCmd("goctl", cmdArgs...)
	fmt.Println("goctl生成模型文件")
	<-c.Start()
	return nil
}

func editLogicFile(name, underlineName, primaryKeyName, primaryKeyJson string, vueFields []map[string]string) error {
	// 新增逻辑
	createLogic := fmt.Sprintf(`
func (l *%vCreateLogic) %vCreate(req *types.%vCreateReq) (err error) {
	var modelParams = new(model.%v)
	err = copier.Copy(modelParams, req)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}
	_, err = l.svcCtx.Feat%vModel.Insert(l.ctx, modelParams)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	return
}
`, name, name, name, name, name)

	// 删除逻辑
	deleteLogic := fmt.Sprintf(`
func (l *%vDeleteLogic) %vDelete(req *types.%vDeleteReq) (err error) {
	err = l.svcCtx.Feat%vModel.Delete(l.ctx, req.%v)
	if err != nil {
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	return
}
`, name, name, name, name, primaryKeyName)

	deletesLogic := fmt.Sprintf(`
func (l *%vDeletesLogic) %vDeletes(req *types.%vDeletesReq) (err error) {
	if len(req.%v) > 0  {
		err = l.svcCtx.Feat%vModel.Deletes(l.ctx, req.%v)
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
	modelParams, err = l.svcCtx.Feat%vModel.FindOne(l.ctx, req.%v)
	if err != nil {
		return errorx2.NewDefaultError(errorx2.UserIdErrorCode)
	}

	err = copier.Copy(modelParams, req)
	if err != nil {
		logx.Error("复制参数失败", err)
		return errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	err = l.svcCtx.Feat%vModel.Update(l.ctx, modelParams)
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
	item, err = l.svcCtx.Feat%vModel.FindOne(l.ctx, req.%v)
	err = copier.Copy(resp, item)
	resp.CreatedAt = utils.Time2Str(item.CreatedAt)
	resp.UpdatedAt = utils.Time2Str(item.UpdatedAt)
	if err != nil {
		logx.Error("复制结果失败", err)
		return nil, err
	}
	return
}
`, name, name, name, name, name, name, name, primaryKeyName)

	// 列表逻辑
	listLogic, _ := getListLogic(name, vueFields)

	// 分页列表逻辑
	pageLogic, _ := getPageLogic(name, vueFields)

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
		if method == "Create" || method == "Update" {
			modifiedContent = strings.Replace(modifiedContent, `"zero-zone/pkg/utils"`, "", -1)
		}
		if method == "Delete" {
			modifiedContent = strings.Replace(modifiedContent, `"github.com/jinzhu/copier"`, "", -1)
			modifiedContent = strings.Replace(modifiedContent, `"zero-zone/applet/model"`, "", -1)
			modifiedContent = strings.Replace(modifiedContent, `"zero-zone/pkg/utils"`, "", -1)
		}
		if method == "Deletes" {
			modifiedContent = strings.Replace(modifiedContent, `"github.com/jinzhu/copier"`, "", -1)
			modifiedContent = strings.Replace(modifiedContent, `"zero-zone/applet/model"`, "", -1)
			modifiedContent = strings.Replace(modifiedContent, `"zero-zone/pkg/utils"`, "", -1)
		}
		if method == "Detail" {
			modifiedContent = strings.Replace(modifiedContent, `errorx2 "zero-zone/pkg/errorx"`, "", -1)
		}
		if method == "List" {
			modifiedContent = strings.Replace(modifiedContent, `"zero-zone/applet/model"`, "", -1)
		}
		if method == "Page" {
			modifiedContent = strings.Replace(modifiedContent, `"zero-zone/applet/model"`, "", -1)
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

func getListLogic(name string, vueFields []map[string]string) (string, error) {
	projectWd, _ := os.Getwd()
	file := filepath.Join(projectWd, "../tpl/list.tpl")
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

func getPageLogic(name string, vueFields []map[string]string) (string, error) {
	projectWd, _ := os.Getwd()
	file := filepath.Join(projectWd, "../tpl/page.tpl")
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

func genWebApiFile(underlineName, lowerCaseName, primaryKeyJson string) error {
	projectWd, _ := os.Getwd()
	// ../../ cmd
	fileDir := filepath.Join(projectWd, "../../../web")
	file := filepath.Join(projectWd, "../tpl/api.tpl")
	fmt.Println("api tpl file", file)
	tpl, err := template.ParseFiles(file)

	if err != nil {
		fmt.Println("create template failed, err:", err)
		return err
	}
	apiFile, err := os.Create(fileDir + "/src/api/feat/" + underlineName + ".js")
	fmt.Println("api file", fileDir+"/src/api/feat/"+underlineName+".js")
	defer apiFile.Close()

	data := map[string]interface{}{
		"UnderlineName":  underlineName,
		"PrimaryKeyJson": primaryKeyJson,
		"LowerCaseName":  lowerCaseName,
	}

	err = tpl.Execute(apiFile, data)
	if err != nil {
		return err
	}
	return nil
}

func genWebVueFile(name, underlineName, primaryKeyJson string, vueFields interface{}, minVueFields interface{}) error {
	projectWd, _ := os.Getwd()
	fileDir := filepath.Join(projectWd, "../../../web")
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
		"MinVueFields":   minVueFields,
	}
	err = tpl.Execute(file, data)
	if err != nil {
		return err
	}
	return nil
}

func getCommentFromGormTag(tag string) string {
	parts := strings.Split(tag, ";") // 根据分号分割tag
	for _, part := range parts {
		//fmt.Println(1, part)
		if strings.Contains(part, "comment:") {
			column := strings.TrimPrefix(part, "comment:")
			//fmt.Println(2, column)
			return column
		}
	}
	return ""
}

func getColumnFromGormTag(tag string) string {
	parts := strings.Split(tag, ";") // 根据分号分割tag
	for _, part := range parts {
		fmt.Println(1, part)
		if strings.Contains(part, "column:") {
			column := strings.TrimPrefix(part, "column:")
			fmt.Println(2, column)
			return column
		}
	}
	return ""
}

// 从大驼峰，转下划线
func GetUnderlineWord(str string) (resp string) {
	var matchNonAlphaNumeric = regexp.MustCompile(`[^a-zA-Z0-9]+`)
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
	str = matchNonAlphaNumeric.ReplaceAllString(str, "_")     //非常规字符转化为 _
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}") //拆分出连续大写
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")  //拆分单词
	return strings.ToLower(snake)
}
