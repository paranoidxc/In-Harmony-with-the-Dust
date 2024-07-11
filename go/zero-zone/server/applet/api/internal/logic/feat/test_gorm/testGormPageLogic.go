package test_gorm

import (
	"context"
	"fmt"
	"strings"

	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"

	errorx2 "zero-zone/pkg/errorx"
)

type TestGormPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTestGormPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TestGormPageLogic {
	return &TestGormPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TestGormPageLogic) TestGormPage(req *types.TestGormPageReq) (resp *types.TestGormPageResp, err error) {
	where := " 1 "
	if len(strings.TrimSpace(req.Text)) > 0 {
		where = where + fmt.Sprintf(" AND text LIKE '%s'", "%"+strings.TrimSpace(req.Text)+"%")
	}

	featTestGormPage, err := l.svcCtx.FeatTestGormModel.FindPageByWhere(where, req.Page, req.Limit)
	if err != nil {
		return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	var item types.TestGorm
	TestGormPage := make([]types.TestGorm, 0)
	for _, v := range featTestGormPage {
		err := copier.Copy(&item, &v)
		if err != nil {
			return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}
		TestGormPage = append(TestGormPage, item)
	}

	total, err := l.svcCtx.FeatTestGormModel.FindPageByWhereCount(where)
	if err != nil {
		return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	pagination := types.Pagination{
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}

	return &types.TestGormPageResp{
		List:       TestGormPage,
		Pagination: pagination,
	}, nil
}
