package logic

import (
	"KyokaSuigetsu/internal/svc"
	"KyokaSuigetsu/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecoveryByTimeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecoveryByTimeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecoveryByTimeLogic {
	return &RecoveryByTimeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecoveryByTimeLogic) RecoveryByTime(req *types.RecoveryByTimeRequest) (resp *types.RecoveryByTimeResponse, err error) {
	r := l.svcCtx.Recovery
	fakeMaster, err := r.RegisterNewFakeMaster(&req.Info)
	if err != nil {
		return
	}
	r.UntilTimestamp(fakeMaster, req.RecoverTimestamp)
	err = r.Start(fakeMaster)
	if err != nil {
		return
	}
	resp = &types.RecoveryByTimeResponse{
		Code:           0,
		Message:        "Success!",
		FakeMasterInfo: fakeMaster,
	}
	return
}
