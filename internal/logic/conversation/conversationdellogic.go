package conversation

import (
	"context"

	"go-chat-sse/internal/svc"
	"go-chat-sse/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConversationdelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConversationdelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConversationdelLogic {
	return &ConversationdelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConversationdelLogic) Conversationdel(req *types.ConversationDelReq) error {
	// todo: add your logic here and delete this line

	return nil
}
