package handler

import (
	"net/http"

	"KyokaSuigetsu/internal/logic"
	"KyokaSuigetsu/internal/svc"
	"KyokaSuigetsu/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func recoveryByGTIDHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RecoveryByGTIDRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewRecoveryByGTIDLogic(r.Context(), svcCtx)
		resp, err := l.RecoveryByGTID(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
