package handler

import (
	"net/http"

	"KyokaSuigetsu/internal/logic"
	"KyokaSuigetsu/internal/svc"
	"KyokaSuigetsu/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func recoveryByTimeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RecoveryByTimeRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewRecoveryByTimeLogic(r.Context(), svcCtx)
		resp, err := l.RecoveryByTime(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
