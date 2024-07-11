package login

import (
	"context"
	errorx2 "zero-zone/pkg/errorx"
	"zero-zone/pkg/globalkey"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"
)

type GetLogLoginPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLogLoginPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLogLoginPageLogic {
	return &GetLogLoginPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLogLoginPageLogic) GetLogLoginPage(req *types.LogLoginPageReq) (resp *types.LogLoginPageResp, err error) {
	loginLogList, err := l.svcCtx.SysLogModel.FindPage(l.ctx, globalkey.SysLoginLogType, req.Page, req.Limit)
	if err != nil {
		return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	var loginLog types.LogLogin
	logList := make([]types.LogLogin, 0)
	for _, v := range loginLogList {
		err := copier.Copy(&loginLog, &v)
		loginLog.CreateTime = v.CreateTime.Format(globalkey.SysDateFormat)
		if err != nil {
			return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}
		logList = append(logList, loginLog)
	}

	total, err := l.svcCtx.SysLogModel.FindCount(l.ctx, globalkey.SysLoginLogType)
	if err != nil {
		return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	pagination := types.Pagination{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}

	return &types.LogLoginPageResp{
		List:       logList,
		Pagination: pagination,
	}, nil
}
