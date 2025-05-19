package chat

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/rest/httpx"
	"go-chat-sse/internal/logic/chat"
	"go-chat-sse/internal/svc"
	"go-chat-sse/internal/types"
	"log"
	"net/http"
)

func MessagePullHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 设置 SSE 相关头信息
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		messageChan := make(chan string)

		go func(ctx context.Context, messageChan chan string) {
			var req types.MessagePullReq
			if err := httpx.Parse(r, &req); err != nil {
				httpx.ErrorCtx(ctx, w, err)
				return
			}
			l := chat.NewMessagePullLogic(ctx, svcCtx)
			l.MessagePullLogic(&req, w, r, messageChan)
			defer close(messageChan)
		}(context.WithoutCancel(r.Context()), messageChan)
		// 创建一个通道用于发送消息

		// 监听客户端断开连接
		//notify := w.(http.CloseNotifier).CloseNotify()

		for {
			select {
			//case <-notify:
			//	log.Println("客户端断开连接")
			//	return
			case msg := <-messageChan:
				// 按照 SSE 格式发送消息
				// data: 开头，结尾两个换行
				fmt.Fprintf(w, "data: %s\n\n", msg)

				// 刷新缓冲区，确保消息立即发送
				if flusher, ok := w.(http.Flusher); ok {
					flusher.Flush()
				} else {
					log.Println("无法初始化刷新器")
				}
			}
		}

	}

}
