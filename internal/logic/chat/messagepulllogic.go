package chat

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
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
	ctx               context.Context
	svcCtx            *svc.ServiceContext
	MessageModel      model.MessageModel
	SessionModel      model.SessionModel
	ConversationModel model.ConversationModel
}

func NewMessagePullLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessagePullLogic {
	return &MessagePullLogic{
		Logger:            logx.WithContext(ctx),
		ctx:               ctx,
		svcCtx:            svcCtx,
		MessageModel:      model.NewMessageModel(svcCtx.Mysql),
		SessionModel:      model.NewSessionModel(svcCtx.Mysql),
		ConversationModel: model.NewConversationModel(svcCtx.Mysql),
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
	var messageReq []third.Message
	// 查询sessionid对应的message
	sessionOne, err := l.SessionModel.FindOneBySessionId(l.ctx, req.SessionId)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			l.Logger.Infof("%s sessionid:%d not found", MessagePull, req.SessionId)
			return biz.MessageNotExistErr
		} else {
			l.Logger.Errorf("%s db error,err:%s", MessagePull, err)
			return biz.DBError
		}
	}
	messageOne, err := l.MessageModel.FindOneByMessageId(l.ctx, sessionOne.MessageId)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			l.Logger.Infof("%s message:%d not found", MessagePull, sessionOne.MessageId)
			return biz.MessageNotExistErr
		} else {
			l.Logger.Errorf("%s db error,err:%s", MessagePull, err)
			return biz.DBError
		}
	}
	//  获取conversation，判断是否为新一轮对话，如果是新一轮对话，要在conversation表新开一个，如果是旧的，那么读取conversation的所有消息记录。
	conversationOne, err := l.ConversationModel.FindOneBySessionId(l.ctx, sessionOne.MessageId)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			// 会话不存在，新建一个conversation
		} else {
			l.Logger.Errorf("%s db error,err:%s", MessagePull, err)
			return biz.DBError
		}
	}
	if conversationOne.MessageNums > 1 {
		// 多轮会话,读取所有的历史消息,利用conversation去message表读取历史消息
		messageMore, err := l.MessageModel.FindMoreByConversationId(l.ctx, conversationOne.Id)
		if err != nil {
			l.Logger.Infof("%s db error,err:%s", MessagePull, err)
			return biz.DBError
		}
		var role string
		for _, message := range *messageMore {
			if message.FromId == 100 {
				role = "system"
			} else {
				role = "user"
			}
			messageReq = append(messageReq, third.Message{
				Role:    role,
				Content: message.Content,
			})
		}
	} else {
		// 单论对话，第一条消息
		messageReq = []third.Message{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: messageOne.Content},
		}
	}

	// 获取ds消息
	requestData := third.ChatRequest{
		Messages: messageReq,
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
