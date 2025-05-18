package user

import (
	"go-chat-sse/internal/biz"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-chat-sse/internal/logic/user"
	"go-chat-sse/internal/svc"
	"go-chat-sse/internal/types"
)

func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, biz.Success(resp))
		}
	}
}
