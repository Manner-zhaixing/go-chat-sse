package chat

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"go-chat-sse/internal/biz"
	"go-chat-sse/internal/model"
	"go-chat-sse/internal/svc"
	"go-chat-sse/internal/third"
	"go-chat-sse/internal/types"
	"io"
	"net/http"
	"strings"
)

const MessagePull = "[messagePull]"

type MessagePullLogic struct {
	logx.Logger
	ctx          context.Context
	svcCtx       *svc.ServiceContext
	MessageModel model.MessageModel
	SessionModel model.SessionModel
}

func NewMessagePullLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessagePullLogic {
	return &MessagePullLogic{
		Logger:       logx.WithContext(ctx),
		ctx:          ctx,
		svcCtx:       svcCtx,
		MessageModel: model.NewMessageModel(svcCtx.Mysql),
		SessionModel: model.NewSessionModel(svcCtx.Mysql),
	}
}

func (l *MessagePullLogic) checkData(req *types.MessagePullReq) error {
	if req.SessionId <= 0 {
		return biz.ParamError
	}
	return nil
}

func (l *MessagePullLogic) MessagePullLogic(req *types.MessagePullReq, w http.ResponseWriter, r *http.Request) error {
	// 1.校验数据
	err := l.checkData(req)
	if err != nil {
		l.Logger.Infof("%s req data param error,err:%s", MessagePull, err)
		return biz.ParamError
	}

	//// 查询sessionid对应的message
	//sessionOne, err := l.SessionModel.FindOneBySessionId(l.ctx, req.SessionId)
	//if err != nil {
	//	if errors.Is(err, sqlx.ErrNotFound) {
	//		l.Logger.Infof("%s sessionid:%d not found", MessagePull, req.SessionId)
	//		return biz.MessageNotExistErr
	//	} else {
	//		l.Logger.Errorf("%s db error,err:%s", MessagePull, err)
	//		return biz.DBError
	//	}
	//}
	//messageOne, err := l.MessageModel.FindOneByMessageId(l.ctx, sessionOne.MessageId)
	//if err != nil {
	//	if errors.Is(err, sqlx.ErrNotFound) {
	//		l.Logger.Infof("%s message:%d not found", MessagePull, sessionOne.MessageId)
	//		return biz.MessageNotExistErr
	//	} else {
	//		l.Logger.Errorf("%s db error,err:%s", MessagePull, err)
	//		return biz.DBError
	//	}
	//}

	// 获取ds消息
	//fmt.Println(messageOne.Content)
	requestData := third.ChatRequest{
		Messages: []third.Message{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: "介绍一下你自己，最多50个字"},
		},
	}

	resp, err := third.StreamChatRequest(requestData)
	if err != nil {
		l.Logger.Infof("%s DeepSeek error,err:%s", MessagePull, err)
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		l.Logger.Infof("%s DeepSeek error,err:%s,body:%s", MessagePull, err, string(body))
		return nil
	}
	// 5. 流式读取DeepSeek响应并转发给客户端
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		// SSE 数据行以 "data: " 开头
		if bytes.HasPrefix(line, []byte("data: ")) {
			data := line[6:] // 去掉 "data: " 前缀
			if strings.Contains(string(data), "[DONE]") {
				l.Logger.Infof("%s stream end,sessionid:%s", MessagePull, req.SessionId)
				break
			}

			// 解析JSON数据
			var chunk third.StreamChatResponse
			if err := json.Unmarshal(data, &chunk); err != nil {
				return fmt.Errorf("解析流数据失败: %w", err)
			}
			// 打印内容
			for _, choice := range chunk.Choices {
				if choice.Delta.Content != "" {
					// fmt.Println(choice.Delta.Content)
					fmt.Fprintf(w, "data: %s\n\n", choice.Delta.Content)
					//fmt.Fprintf(w, "data: %s\n\n", choice)
					if flusher, ok := w.(http.Flusher); ok {
						flusher.Flush()
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取流数据失败: %w", err)
	}
	return nil
}
