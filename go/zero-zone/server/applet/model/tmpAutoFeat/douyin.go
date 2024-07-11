package tmpAutoFeat

// gorm:"column:firm_name;not null;uniqueIndex;type:varchar(255);comment:厂商名称;"`
// gorm tag 暂时只用到 column 和 comment

type TmpDouyin struct {
	Id           int64  `json:"id" gorm:"column:id;comment:唯一id;"`
	Name         string `json:"name" label:"应用名称" gorm:"column:name;comment:应用名称;"`
	AppId        string `json:"appId" label:"AppID" gorm:"column:app_id;comment:AppID;"`
	AppSecret    string `json:"appSecret" label:"AppSecret" gorm:"column:app_secret;comment:AppSecret;"`
	ClientToken  string `json:"clientToken" label:"client_token" gorm:"column:app_client_tokensecret;comment:client_token;"`
	Xcode        string `json:"xcode" label:"临时调用凭证" gorm:"column:xcode;comment:临时调用凭证;"`
	AccessToken  string `json:"accessToken" label:"接口调用凭证" gorm:"column:access_token;comment:接口调用凭证;"`
	RefreshToken string `json:"refreshToken" label:"刷新令牌" gorm:"column:refresh_token;comment:刷新令牌;"`
	CreatedAt    string `json:"createdAt" gorm:"column:created_at;comment:创建时间;"`
	UpdatedAt    string `json:"updatedAt" gorm:"column:updated_at;comment:更新时间;"`
	DeletedAt    string `json:"deletedAt" gorm:"column:deleted_at;comment:删除时间;"`
}

type TmpThirdPartDevConf struct {
	Id        int64  `json:"id" gorm:"column:id;comment:唯一id;"`
	Name      string `json:"name" label:"应用名称" gorm:"column:name;comment:应用名称;"`
	AppId     string `json:"appId" label:"AppID" gorm:"column:app_id;comment:AppID;"`
	AppSecret string `json:"appSecret" label:"AppSecret" gorm:"column:app_secret;comment:AppSecret;"`
	CreatedAt string `json:"createdAt" gorm:"column:created_at;comment:创建时间;"`
	UpdatedAt string `json:"updatedAt" gorm:"column:updated_at;comment:更新时间;"`
	DeletedAt string `json:"deletedAt" gorm:"column:deleted_at;comment:删除时间;"`
	Typo      int    `json:"typo" gorm:"column:typo;comment:渠道;"`
}

type TmpSaasCooperateAuth struct {
	Id        int64  `json:"id" gorm:"column:id;comment:唯一id;"`
	DevConfId int64  `json:"devConfId" label:"DEV应用ID" gorm:"column:dev_conf_ide;comment:DEV应用ID;"`
	Content   string `json:"content" label:"授权信息" gorm:"column:content;comment:授权信息;"`
	CreatedAt string `json:"createdAt" gorm:"column:created_at;comment:创建时间;"`
	UpdatedAt string `json:"updatedAt" gorm:"column:updated_at;comment:更新时间;"`
	DeletedAt string `json:"deletedAt" gorm:"column:deleted_at;comment:删除时间;"`
}

type TmpCooperateShop struct {
	Id                  int64  `json:"id" gorm:"column:id;comment:唯一id;"`
	Account             string `json:"account" label:"接口账号" gorm:"column:account;comment:接口账号;"`
	ApiKey              string `json:"apiKey" label:"接口key" gorm:"column:api_key;comment:接口key;"`
	Name                string `json:"name" label:"门店信息" gorm:"column:content;comment:门店信息;"`
	Typo                int    `json:"typo" label:"渠道" gorm:"column:typo;comment:渠道;"`
	SaasCooperateAuthId int64  `json:"saasCooperateAuthId" label:"渠道商户" gorm:"column:saas_cooperate_auth_id;comment:渠道商户;"`
	ShopId              string `json:"shopId" label:"渠道商户门店" gorm:"column:shop_id;comment:渠道商户门店;"`
	ShopName            string `json:"shopName" label:"渠道商户门店" gorm:"column:shop_name;comment:渠道商户门店;"`
	CreatedAt           string `json:"createdAt" gorm:"column:created_at;comment:创建时间;"`
	UpdatedAt           string `json:"updatedAt" gorm:"column:updated_at;comment:更新时间;"`
	DeletedAt           string `json:"deletedAt" gorm:"column:deleted_at;comment:删除时间;"`
}
