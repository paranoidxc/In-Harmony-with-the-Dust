package types

type LoginReq struct {
	CaptchaId  string `json:"captchaId"   label:"验证码id"`
	VerifyCode string `json:"verifyCode"  label:"验证码"`
	Account    string `json:"account"     label:"账号"`
	Username   string `json:"username"    label:"账号"`
	Password   string `json:"password"    label:"密码"`
}

type LoginResp struct {
	Token      string `json:"token"`
	TokenName  string `json:"tokenName"`
	TokenValue string `json:"tokenValue"`
}

type UserInfoResp struct {
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

type UserProfileInfoResp struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Gender   int64  `json:"gender"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Remark   string `json:"remark"`
	Avatar   string `json:"avatar"`
}

type UpdateProfileReq struct {
	Username string `json:"username"  validate:"required,min=2,max=12"   label:"姓名"`
	Nickname string `json:"nickname"  validate:"omitempty,min=2,max=12"  label:"昵称"`
	Gender   int64  `json:"gender"    validate:"gte=0,lte=2"             label:"性别"`
	Email    string `json:"email"     validate:"omitempty,email"         label:"邮箱"`
	Mobile   string `json:"mobile"    validate:"omitempty,len=11"        label:"手机号"`
	Avatar   string `json:"avatar"    validate:"required,url"            label:"头像"`
}

type Menu struct {
	Id           int64  `json:"id"`
	ParentId     int64  `json:"parentId"`
	Name         string `json:"name"`
	Router       string `json:"router"`
	Type         int64  `json:"type"`
	Icon         string `json:"icon"`
	OrderNum     int64  `json:"orderNum"`
	ViewPath     string `json:"viewPath"`
	IsShow       int64  `json:"isShow"`
	ActiveRouter string `json:"activeRouter"`
}

type NewMenu struct {
	NewTest          string    `json:"newTest"`
	Component        string    `json:"component"`
	Icon             string    `json:"icon"`
	IsHasPerm        bool      `json:"isHasPerm"`
	IsShow           bool      `json:"isShow"`
	KeepAlive        bool      `json:"keepAlive"`
	IsShowBreadcrumb bool      `json:"isShowBreadcrumb"`
	MenuId           int64     `json:"menuId"`
	ParentId         int64     `json:"parentId"`
	ParentName       string    `json:"parentName"`
	Path             string    `json:"path"`
	PermList         string    `json:"permList"`
	redirect         string    `json:"redirect"`
	sort             int       `json:"sort"`
	Title            string    `json:"title"`
	OrderNum         int64     `json:"orderNum"`
	Children         []NewMenu `json:"children"`
}

type UserPermMenuResp struct {
	Menus              []Menu    `json:"menus"`
	PermissionTreeList []NewMenu `json:"permissionTreeList"`
	Perms              []string  `json:"perms"`
}

type UpdatePasswordReq struct {
	OldPassword string `json:"oldPassword"  validate:"min=6,max=120"  label:"旧密码"`
	NewPassword string `json:"newPassword"  validate:"min=6,max=120"  label:"新密码"`
}

type LoginCaptchaResp struct {
	CaptchaId  string `json:"captchaId"`
	VerifyCode string `json:"verifyCode"`
}

type GenerateAvatarResp struct {
	AvatarUrl string `json:"avatarUrl"`
}

type PermMenu struct {
	Id           int64      `json:"id"`
	Value        int64      `json:"value"`
	ParentId     int64      `json:"parentId"`
	Name         string     `json:"name"`
	Label        string     `json:"label"`
	Router       string     `json:"router"`
	Perms        []string   `json:"perms"`
	Type         int64      `json:"type"`
	Icon         string     `json:"icon"`
	OrderNum     int64      `json:"orderNum"`
	ViewPath     string     `json:"viewPath"`
	IsShow       int64      `json:"isShow"`
	KeepAlive    int        `json:"keepAlive"`
	ActiveRouter string     `json:"activeRouter"`
	Children     []PermMenu `json:"children"`
}

type SysPermMenuListResp struct {
	List []PermMenu `json:"list"`
}

type AddSysPermMenuReq struct {
	ParentId     int64    `json:"parentId"      validate:"number,gte=0"           label:"父级菜单id"`
	Name         string   `json:"name"          validate:"min=2,max=50"           label:"菜单名称"`
	Router       string   `json:"router"        validate:"omitempty,max=1024"     label:"路由"`
	Perms        []string `json:"perms"         validate:"omitempty,unique"       label:"权限"`
	Type         int64    `json:"type"          validate:"number,gte=0,lte=2"     label:"类型"`
	Icon         string   `json:"icon"          validate:"omitempty,max=200"      label:"图标"`
	OrderNum     int64    `json:"orderNum"      validate:"number,gte=0,lte=9999"  label:"排序"`
	ViewPath     string   `json:"viewPath"      validate:"omitempty,max=1024"     label:"视图路径"`
	IsShow       int64    `json:"isShow"        validate:"number,gte=0,lte=1"     label:"显示状态"`
	KeepAlive    int64    `json:"keepAlive"     validate:"number,gte=0,lte=1"     label:"KeepAlive"`
	ActiveRouter string   `json:"activeRouter"  validate:"omitempty,max=1024"     label:"激活路由"`
}

type DeleteSysPermMenuReq struct {
	Id int64 `json:"id"  validate:"number,gte=1" label:"菜单id"`
}

type UpdateSysPermMenuReq struct {
	Id           int64    `json:"id"            validate:"number,gte=1"           label:"菜单id"`
	ParentId     int64    `json:"parentId"      validate:"number,gte=0"           label:"父级菜单id"`
	Name         string   `json:"name"          validate:"min=2,max=50"           label:"菜单名称"`
	Router       string   `json:"router"        validate:"omitempty,max=1024"     label:"路由"`
	Perms        []string `json:"perms"         validate:"omitempty,unique"       label:"权限"`
	Type         int64    `json:"type"          validate:"number,gte=0,lte=2"     label:"类型"`
	Icon         string   `json:"icon"          validate:"omitempty,max=200"      label:"图标"`
	OrderNum     int64    `json:"orderNum"      validate:"number,gte=0,lte=9999"  label:"排序"`
	ViewPath     string   `json:"viewPath"      validate:"omitempty,max=1024"     label:"视图路径"`
	IsShow       int64    `json:"isShow"        validate:"number,gte=0,lte=1"     label:"显示状态"`
	KeepAlive    int64    `json:"keepAlive"     validate:"number,gte=0,lte=1"     label:"KeepAlive"`
	ActiveRouter string   `json:"activeRouter"  validate:"omitempty,max=1024"     label:"激活路由"`
}

type Role struct {
	Id          int64   `json:"id"`
	ParentId    int64   `json:"parentId"`
	Name        string  `json:"name"`
	UniqueKey   string  `json:"uniqueKey"`
	PermMenuIds []int64 `json:"permMenuIds"`
	Remark      string  `json:"remark"`
	Status      int64   `json:"status"`
	OrderNum    int64   `json:"orderNum"`
}

type SysRoleListResp struct {
	List []Role `json:"list"`
}

type AddSysRoleReq struct {
	ParentId       int64   `json:"parentId"     validate:"number,gte=0"          label:"父级角色id"`
	Name           string  `json:"name"         validate:"min=2,max=50"          label:"角色名称"`
	UniqueKey      string  `json:"uniqueKey"    validate:"min=2,max=50"          label:"角色标识"`
	PermMenuIds    []int64 `json:"permMenuIds"  validate:"omitempty,unique"      label:"权限ids"`
	Remark         string  `json:"remark,optional"                     label:"备注"`
	Status         int64   `json:"status"       validate:"number,gte=0,lte=1"    label:"状态"`
	OrderNum       int64   `json:"orderNum"     validate:"number,gte=0,lte=9999" label:"排序"`
	PermMenuIdsAll []int64 `json:"permMenuIdsAll"  validate:"omitempty"      label:"权限ids"`
}

type DeleteSysRoleReq struct {
	Id int64 `json:"id"  validate:"number,gte=2" label:"角色id"`
}

type UpdateSysRoleReq struct {
	Id             int64   `json:"id"           validate:"number,gte=1"           label:"角色id"`
	ParentId       int64   `json:"parentId"     validate:"number,gte=0"           label:"父级角色id"`
	Name           string  `json:"name"         validate:"min=2,max=50"           label:"角色名称"`
	UniqueKey      string  `json:"uniqueKey"    validate:"min=2,max=50"           label:"角色标识"`
	PermMenuIds    []int64 `json:"permMenuIds"  validate:"omitempty,unique"       label:"权限ids"`
	Remark         string  `json:"remark,optional"                      label:"备注"`
	Status         int64   `json:"status"       validate:"number,gte=0,lte=1"     label:"状态"`
	OrderNum       int64   `json:"orderNum"     validate:"number,gte=0,lte=9999"  label:"排序"`
	PermMenuIdsAll []int64 `json:"permMenuIdsAll"  validate:"omitempty"      label:"权限ids"`
}

type Dept struct {
	Id        int64  `json:"id"`
	ParentId  int64  `json:"parentId"`
	Name      string `json:"name"`
	FullName  string `json:"fullName"`
	UniqueKey string `json:"uniqueKey"`
	Type      int64  `json:"type"`
	Status    int64  `json:"status"`
	OrderNum  int64  `json:"orderNum"`
	Remark    string `json:"remark"`
}

type SysDeptPageReq struct {
	PageReq
}

type SysDeptPageResp struct {
	List []Dept `json:"list"`
	Pagination
}

type SysDeptListResp struct {
	List []Dept `json:"list"`
}

type AddSysDeptReq struct {
	ParentId  int64  `json:"parentId"   validate:"number,gte=0"            label:"父级部门id"`
	Name      string `json:"name"       validate:"min=2,max=50"            label:"部门名称"`
	FullName  string `json:"fullName"   validate:"omitempty,min=2,max=50"  label:"部门全称"`
	UniqueKey string `json:"uniqueKey"  validate:"min=2,max=50"            label:"部门标识"`
	Type      int64  `json:"type"       validate:"number,gte=1,lte=3"      label:"部门类型"`
	Status    int64  `json:"status"     validate:"number,gte=0,lte=1"      label:"状态"`
	OrderNum  int64  `json:"orderNum"   validate:"number,gte=0,lte=9999"   label:"排序"`
	Remark    string `json:"remark"     validate:"max=200"                 label:"备注"`
}

type SearchDeptReq struct {
	Name     string `form:"name,optional"`
	FullName string `form:"fullName,optional"`
	Status   int    `form:"status,optional"`
	PageReq
}

type DeleteSysDeptReq struct {
	Id int64 `json:"id"  validate:"number,gte=1" label:"部门id"`
}

type UpdateSysDeptReq struct {
	Id        int64  `json:"id"         validate:"number,gte=1"            label:"部门id"`
	ParentId  int64  `json:"parentId"   validate:"number,gte=0"            label:"父级部门id"`
	Name      string `json:"name"       validate:"min=2,max=50"            label:"部门名称"`
	FullName  string `json:"fullName"   validate:"omitempty,min=2,max=50"  label:"部门全称"`
	UniqueKey string `json:"uniqueKey"  validate:"min=2,max=50"            label:"部门标识"`
	Type      int64  `json:"type"       validate:"number,gte=1,lte=3"      label:"部门类型"`
	Status    int64  `json:"status"     validate:"number,gte=0,lte=1"      label:"状态"`
	OrderNum  int64  `json:"orderNum"   validate:"number,gte=0,lte=9999"   label:"排序"`
	Remark    string `json:"remark"     validate:"max=200"                 label:"备注"`
}

type Job struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Status   int64  `json:"status"`
	OrderNum int64  `json:"orderNum"`
}

type SysJobPageReq struct {
	PageReq
}

type SysJobPageResp struct {
	List       []Job      `json:"list"`
	Pagination Pagination `json:"pagination"`
}

type AddSysJobReq struct {
	Name     string `json:"name"      validate:"min=2,max=50"           label:"岗位名称"`
	Status   int64  `json:"status"    validate:"number,gte=0,lte=1"     label:"状态"`
	OrderNum int64  `json:"orderNum"  validate:"number,gte=0,lte=9999"  label:"排序"`
}

type DeleteSysJobReq struct {
	Id int64 `json:"id"  validate:"number,gte=1" label:"岗位id"`
}

type UpdateSysJobReq struct {
	Id       int64  `json:"id"        validate:"number,gte=1"           label:"岗位id"`
	Name     string `json:"name"      validate:"min=2,max=50"           label:"岗位名称"`
	Status   int64  `json:"status"    validate:"number,gte=0,lte=1"     label:"状态"`
	OrderNum int64  `json:"orderNum"  validate:"number,gte=0,lte=9999"  label:"排序"`
}

type Profession struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Status   int64  `json:"status"`
	OrderNum int64  `json:"orderNum"`
}

type SysProfessionPageReq struct {
	PageReq
}

type SysProfessionPageResp struct {
	List       []Profession `json:"list"`
	Pagination Pagination   `json:"pagination"`
}

type AddSysProfessionReq struct {
	Name     string `json:"name"      validate:"min=2,max=50"           label:"职称"`
	Status   int64  `json:"status"    validate:"number,gte=0,lte=1"     label:"状态"`
	OrderNum int64  `json:"orderNum"  validate:"number,gte=0,lte=9999"  label:"排序"`
}

type DeleteSysProfessionReq struct {
	Id int64 `json:"id"  validate:"number,gte=1" label:"职称id"`
}

type UpdateSysProfessionReq struct {
	Id       int64  `json:"id"        validate:"number,gte=1"           label:"职称id"`
	Name     string `json:"name"      validate:"min=2,max=50"           label:"职称"`
	Status   int64  `json:"status"    validate:"number,gte=0,lte=1"     label:"状态"`
	OrderNum int64  `json:"orderNum"  validate:"number,gte=0,lte=9999"  label:"排序"`
}

type UserProfession struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type UserJob struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type UserDept struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type UserRole struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type User struct {
	Id       int64  `json:"id"`
	Account  string `json:"account"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Gender   int64  `json:"gender"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	//Profession UserProfession `json:"profession"`
	//Job        UserJob        `json:"job"`
	//Dept       UserDept       `json:"dept"`
	Roles    []UserRole `json:"roles"`
	Status   int64      `json:"status"`
	OrderNum int64      `json:"orderNum"`
	Remark   string     `json:"remark"`
}

type SysUserPageReq struct {
	PageReq
	DeptId int64 `form:"deptId"  validate:"number,gte=0" label:"部门id"`
}

type SysUserPageResp struct {
	List       []User     `json:"list"`
	Pagination Pagination `json:"pagination"`
}

type DetailSysUserReq struct {
	Id int64 `form:"id"  validate:"number,gte=2" label:"用户id"`
}

type DetailSysUserResp struct {
	User
}

type AddSysUserReq struct {
	Account      string  `json:"account,optional"       validate:"required,min=4,max=50"            label:"账号"`
	Username     string  `json:"username,optional"      validate:"required,min=2,max=50"            label:"姓名"`
	Nickname     string  `json:"nickname,optional"      validate:"omitempty,min=2,max=50"  label:"昵称"`
	Gender       int64   `json:"gender,optional"        validate:"number,gte=0,lte=2"      label:"性别"`
	Email        string  `json:"email,optional"         validate:"omitempty,email"         label:"邮箱"`
	Mobile       string  `json:"mobile,optional"        validate:"omitempty,min=11"        label:"手机号"`
	ProfessionId int64   `json:"professionId,optional"              label:"职称id"`
	JobId        int64   `json:"jobId,optional"                    label:"岗位id"`
	DeptId       int64   `json:"deptId,optional"                    label:"部门id"`
	RoleIds      []int64 `json:"roleIds,optional"       validate:"unique"                  label:"角色ids"`
	Status       int64   `json:"status"        validate:"number,gte=0,lte=1"      label:"状态"`
	OrderNum     int64   `json:"orderNum,optional"      validate:"number,gte=0,lte=9999"   label:"排序"`
	Remark       string  `json:"remark,optional"        validate:"max=200"                 label:"备注"`
}

type DeleteSysUserReq struct {
	Id int64 `json:"id"  validate:"number,gte=2" label:"用户id"`
}

type UpdateSysUserReq struct {
	Id           int64   `json:"id"            validate:"number,gte=2"            label:"用户id"`
	Username     string  `json:"username"      validate:"min=2,max=50"            label:"姓名"`
	Nickname     string  `json:"nickname"      validate:"omitempty,min=2,max=50"  label:"昵称"`
	Gender       int64   `json:"gender"        validate:"number,gte=0,lte=2"      label:"性别"`
	Email        string  `json:"email"         validate:"omitempty,email"         label:"邮箱"`
	Mobile       string  `json:"mobile"        validate:"omitempty,min=11"        label:"手机号"`
	ProfessionId int64   `json:"professionId,optional"             label:"职称id"`
	JobId        int64   `json:"jobId,optional"                    label:"岗位id"`
	DeptId       int64   `json:"deptId,optional"                   label:"部门id"`
	RoleIds      []int64 `json:"roleIds"       validate:"unique"                  label:"角色ids"`
	Status       int64   `json:"status,optional"        validate:"number,gte=0,lte=1"      label:"状态"`
	OrderNum     int64   `json:"orderNum,optional"      validate:"number,gte=0,lte=9999"   label:"排序"`
	Remark       string  `json:"remark,optional"        validate:"max=200"                 label:"备注"`
}

type UpdateSysUserPasswordReq struct {
	Id       int64  `json:"id"        validate:"number,gte=2"  label:"用户id"`
	Password string `json:"password,optional"  validate:"omitempty,min=6,max=12"  label:"密码"`
}

type GetSysUserRdpjInfoReq struct {
	UserId int64 `form:"userId"  validate:"number,gte=0"  label:"用户id"`
}

type Rdpj struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type DeptTree struct {
	Id       int64  `json:"id"`
	ParentId int64  `json:"parentId"`
	Name     string `json:"name"`
}

type RoleTree struct {
	Id       int64  `json:"id"`
	ParentId int64  `json:"parentId"`
	Name     string `json:"name"`
}

type GetSysUserRdpjInfoResp struct {
	Role       []RoleTree `json:"role"`
	Dept       []DeptTree `json:"dept"`
	Profession []Rdpj     `json:"profession"`
	Job        []Rdpj     `json:"job"`
}

type ConfigDict struct {
	Id        int64  `json:"id"`
	ParentId  int64  `json:"parentId"`
	Name      string `json:"name"`
	Type      int64  `json:"type"`
	UniqueKey string `json:"uniqueKey"`
	Value     string `json:"value"`
	OrderNum  int64  `json:"orderNum"`
	Remark    string `json:"remark"`
	Status    int64  `json:"status"`
}

type ConfigDictListResp struct {
	List []ConfigDict `json:"list"`
}

type ConfigDictPageReq struct {
	PageReq
	ParentId int64 `form:"parentId"  validate:"number,gte=0" label:"字典集id"`
}

type ConfigDictPageResp struct {
	List       []ConfigDict `json:"list"`
	Pagination Pagination   `json:"pagination"`
}

type AddConfigDictReq struct {
	ParentId  int64  `json:"parentId"   validate:"number,gte=0"         label:"字典集id"`
	Name      string `json:"name"       validate:"min=2,max=50"         label:"名称"`
	Type      int64  `json:"type"       validate:"number,gte=1,lte=12"  label:"类型"`
	UniqueKey string `json:"uniqueKey"  validate:"min=2,max=50"         label:"标识"`
	Value     string `json:"value"      validate:"max=2048"             label:"字典项值"`
	OrderNum  int64  `json:"orderNum"   validate:"gte=0,lte=9999"       label:"排序"`
	Remark    string `json:"remark"     validate:"max=200"              label:"备注"`
	Status    int64  `json:"status"     validate:"number,gte=0,lte=1"   label:"状态"`
}

type DeleteConfigDictReq struct {
	Id int64 `json:"id"  validate:"number,gte=1" label:"字典id"`
}

type UpdateConfigDictReq struct {
	Id       int64  `json:"id"         validate:"number,gte=1"         label:"字典id"`
	ParentId int64  `json:"parentId"   validate:"number,gte=0"         label:"字典集id"`
	Name     string `json:"name"       validate:"min=2,max=50"         label:"名称"`
	Type     int64  `json:"type"       validate:"number,gte=1,lte=12"  label:"类型"`
	Value    string `json:"value"      validate:"max=2048"             label:"字典项值"`
	OrderNum int64  `json:"orderNum"   validate:"gte=0,lte=9999"       label:"排序"`
	Remark   string `json:"remark"     validate:"max=200"              label:"备注"`
	Status   int64  `json:"status"     validate:"number,gte=0,lte=1"   label:"状态"`
}

type LogLogin struct {
	Id         int64  `json:"id"`
	Account    string `json:"account"`
	Ip         string `json:"ip"`
	Uri        string `json:"uri"`
	Status     int64  `json:"status"`
	CreateTime string `json:"createTime"`
}

type LogLoginPageReq struct {
	PageReq
}

type LogLoginPageResp struct {
	List       []LogLogin `json:"list"`
	Pagination Pagination `json:"pagination"`
}

type AutoCurd struct {
	Name string `json:"name"`
}

type AutoCurdCreateReq struct {
	Name         string `json:"name" validate:"required" label:"模型名称"`
	IsAll        int    `json:"isAll" label:"是否全部CURD文件"`
	IsApi        int    `json:"isApi"  label:"Api文件"`
	IsHandle     int    `json:"isHandle"  label:"Handle/Logic文件"`
	IsModel      int    `json:"isModel"  label:"Model文件"`
	IsVue        int    `json:"isVue"  label:"Vue文件"`
	IsLogicWrite int    `json:"isLogicWrite"  label:"Logic文件"`
	IsMenu       int    `json:"isMenu"  label:"菜单权限"`
}

type AutoCurdListResp struct {
	List []AutoCurd `json:"list"`
}
