package test_gorm

import (
	"context"

	"zero-zone/applet/api/internal/svc"
	"zero-zone/applet/api/internal/types"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"

	errorx2 "zero-zone/pkg/errorx"
)

type TestGormListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTestGormListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TestGormListLogic {
	return &TestGormListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TestGormListLogic) TestGormList(req *types.TestGormListReq) (resp *types.TestGormListResp, err error) {
	where := " 1 "
	/*
	   if len(strings.TrimSpace(req.Text)) > 0 {
	       where = where + fmt.Sprintf(" AND  LIKE '%s'", "%"+strings.TrimSpace(req.Text)+"%")
	   }
	*/
	featTestGormList, err := l.svcCtx.FeatTestGormModel.FindAllByWhere(where)
	if err != nil {
		return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
	}

	var item types.TestGorm
	TestGormList := make([]types.TestGorm, 0)
	for _, v := range featTestGormList {
		err := copier.Copy(&item, &v)
		if err != nil {
			return nil, errorx2.NewSystemError(errorx2.ServerErrorCode, err.Error())
		}
		TestGormList = append(TestGormList, item)
	}

	return &types.TestGormListResp{
		List: TestGormList,
	}, nil
}
