package user

import (
	"go-chat-sse/internal/biz"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-chat-sse/internal/logic/user"
	"go-chat-sse/internal/svc"
	"go-chat-sse/internal/types"
)

func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewRegisterLogic(r.Context(), svcCtx)
		resp, err := l.Register(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, biz.Success(resp))
		}
	}
}
