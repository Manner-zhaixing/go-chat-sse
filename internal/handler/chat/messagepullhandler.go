package chat

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"go-chat-sse/internal/logic/chat"
	"go-chat-sse/internal/svc"
	"go-chat-sse/internal/types"
	"net/http"
)

func MessagePullHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MessagePullReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		// 设置sse响应头
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		l := chat.NewMessagePullLogic(r.Context(), svcCtx)
		err := l.MessagePullLogic(&req, w, r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		}
	}
}
