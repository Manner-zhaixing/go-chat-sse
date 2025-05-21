package conversation

import (
	"go-chat-sse/internal/biz"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-chat-sse/internal/logic/conversation"
	"go-chat-sse/internal/svc"
)

func ConversationaddHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := conversation.NewConversationaddLogic(r.Context(), svcCtx)
		resp, err := l.Conversationadd()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, biz.Success(resp))
		}
	}
}
