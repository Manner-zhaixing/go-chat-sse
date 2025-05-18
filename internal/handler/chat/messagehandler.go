package chat

import (
	"go-chat-sse/internal/biz"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-chat-sse/internal/logic/chat"
	"go-chat-sse/internal/svc"
	"go-chat-sse/internal/types"
)

func MessageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MessageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := chat.NewMessageLogic(r.Context(), svcCtx)
		resp, err := l.Message(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, biz.Success(resp))
		}
	}
}
