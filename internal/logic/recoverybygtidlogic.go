package logic

import (
	"context"

	"KyokaSuigetsu/internal/svc"
	"KyokaSuigetsu/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecoveryByGTIDLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecoveryByGTIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecoveryByGTIDLogic {
	return &RecoveryByGTIDLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecoveryByGTIDLogic) RecoveryByGTID(req *types.RecoveryByGTIDRequest) (resp *types.RecoveryByGTIDResponse, err error) {
	// todo: add your logic here and delete this line
	return
}
