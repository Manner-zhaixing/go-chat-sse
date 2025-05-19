package chat

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-chat-sse/internal/biz"
	"go-chat-sse/internal/model"
	"go-chat-sse/internal/svc"
	"go-chat-sse/internal/third"
	"go-chat-sse/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

const MessagePull = "[messagePull]"

type MessagePullLogic struct {
	logx.Logger
	ctx          context.Context
	svcCtx       *svc.ServiceContext
	MessageModel model.MessageModel
}

func NewMessagePullLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessagePullLogic {
	return &MessagePullLogic{
		Logger:       logx.WithContext(ctx),
		ctx:          ctx,
		svcCtx:       svcCtx,
		MessageModel: model.NewMessageModel(svcCtx.Mysql),
	}
}

func (l *MessagePullLogic) checkData(req *types.MessagePullReq) error {
	if req.SessionId <= 0 {
		return biz.ParamError
	}
	return nil
}

func (l *MessagePullLogic) MessagePullLogic(req *types.MessagePullReq, w http.ResponseWriter, r *http.Request) (*types.MessagePullResp, error) {
	// 1.校验数据
	err := l.checkData(req)
	if err != nil {
		l.Logger.Infof("%s req data param error,err:%s", MessagePull, err)
		return nil, biz.ParamError
	}
	// 查询sessionid对应的message
	messageOne, err := l.MessageModel.FindOneByMessageId(l.ctx, req.SessionId)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			l.Logger.Infof("%s sessionid:%d not found", MessagePull, req.SessionId)
			return nil, biz.MessageNotExistErr
		} else {
			l.Logger.Errorf("%s db error,err:%s", MessagePull, err)
			return nil, biz.DBError
		}
	}

	// 2.设置响应头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 创建一个通道用于发送消息
	defer close(third.DsDataChannel)

	// 监听客户端断开连接
	notify := w.(http.CloseNotifier).CloseNotify()
	// 获取ds消息
	requestData := third.ChatRequest{
		Messages: []third.Message{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: messageOne.Content},
		},
	}
	err = third.StreamChatRequest(requestData)
	if err != nil {
		l.Logger.Infof("%s deepseek-req error,err:%s", MessagePull, err)
		return nil, biz.ParamError
	}

	for {
		select {
		case <-notify:
			l.Logger.Infof("%s client disconnected", MessagePull)
			return nil, nil
		case msg := <-third.DsDataChannel:
			// 按照 SSE 格式发送消息
			// data: 开头，结尾两个换行
			fmt.Fprintf(w, "data: %s\n\n", msg)

			// 刷新缓冲区，确保消息立即发送
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			} else {
				l.Logger.Infof("%s 无法初始化刷新器.", MessagePull)
			}
		}
	}
}
