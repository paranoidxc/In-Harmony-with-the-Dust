package model

import "zero-zone/applet/model/tmpAutoFeat"

// map key 首字母小写

var AutoCrudModelList = map[string]interface{}{
	"DemoCurd":  tmpAutoFeat.TmpDemoCurd{},
	"SysRegion": tmpAutoFeat.TmpSysRegion{},
	//"ThirdPartDevConf": tmpAutoFeat.TmpThirdPartDevConf{},
	//"SaasCooperateAuth": tmpAutoFeat.TmpSaasCooperateAuth{},
	//"CooperateShop": tmpAutoFeat.TmpCooperateShop{},
	//"TestGorm": tmpAutoFeat.TestGorm{},
	//"HxOrder":  tmpAutoFeat.TmpHxOrder{},
	//"uhxOrder": tmpAutoFeat.TmpUhxOrder{},
}
