package conversation

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-chat-sse/internal/biz"
	"go-chat-sse/internal/model"

	"go-chat-sse/internal/svc"
	"go-chat-sse/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const ConversationmessagesModule = "[ConversationmessagesModule]"

type ConversationmessagesLogic struct {
	logx.Logger
	ctx          context.Context
	svcCtx       *svc.ServiceContext
	MessageModel model.MessageModel
}

func NewConversationmessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConversationmessagesLogic {
	return &ConversationmessagesLogic{
		Logger:       logx.WithContext(ctx),
		ctx:          ctx,
		svcCtx:       svcCtx,
		MessageModel: model.NewMessageModel(svcCtx.Mysql),
	}
}

func (l *ConversationmessagesLogic) Conversationmessages(req *types.ConversationMessageReq) (*types.ConversationMessageResp, error) {
	messageInfos, err := l.MessageModel.FindMoreByConversationId(l.ctx, req.ConversationId)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			l.Logger.Infof("%s not find", ConversationmessagesModule)
			return nil, nil
		} else {
			l.Logger.Errorf("%s db error,err:%s", ConversationmessagesModule)
			return nil, biz.DBError
		}
	}
	var temp []types.ConversationMessage
	for _, messageInfo := range *messageInfos {
		temp = append(temp, types.ConversationMessage{
			Content:        messageInfo.Content,
			ConversationId: messageInfo.ConversationId,
			CurTime:        messageInfo.CurTime.String(),
			Done:           int(messageInfo.Done),
			FromId:         messageInfo.FromId,
			MessageId:      messageInfo.MessageId,
			ModelId:        int(messageInfo.ModelId),
			ToId:           messageInfo.ToId,
		})
	}
	return &types.ConversationMessageResp{
		ConversationMessages: temp,
	}, nil
}
