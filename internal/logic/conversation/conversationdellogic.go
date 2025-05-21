package conversation

import (
	"context"
	"encoding/json"
	"go-chat-sse/internal/biz"
	"go-chat-sse/internal/model"

	"go-chat-sse/internal/svc"
	"go-chat-sse/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const conversationDelModule = "[conversationDelModule]"

type ConversationdelLogic struct {
	logx.Logger
	ctx                   context.Context
	svcCtx                *svc.ServiceContext
	ConversationModel     model.ConversationModel
	UserConversationModel model.UserConversationModel
	UserModel             model.UserModel
	MessageModel          model.MessageModel
}

func NewConversationdelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConversationdelLogic {
	return &ConversationdelLogic{
		Logger:                logx.WithContext(ctx),
		ctx:                   ctx,
		svcCtx:                svcCtx,
		ConversationModel:     model.NewConversationModel(svcCtx.Mysql),
		UserConversationModel: model.NewUserConversationModel(svcCtx.Mysql),
		UserModel:             model.NewUserModel(svcCtx.Mysql),
		MessageModel:          model.NewMessageModel(svcCtx.Mysql),
	}
}

func (l *ConversationdelLogic) Conversationdel(req *types.ConversationDelReq) error {
	// 获取userid
	userid, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		l.Logger.Infof("%s token error", conversationDelModule)
		return biz.TokenError
	}
	// 删除conversation表的信息
	err = l.ConversationModel.Delete(l.ctx, req.ConversationId)
	if err != nil {
		l.Logger.Errorf("%s delete conversation error,DB error,info:%s", conversationDelModule, req)
		return biz.DBError
	}
	// 删除user_conversation表的信息
	err = l.UserConversationModel.DeleteByConversationIdAndUserId(l.ctx, req.ConversationId, userid)
	if err != nil {
		l.Logger.Errorf("%s delete user_conversation error,DB error,info:%s", conversationDelModule, req)
		return biz.DBError
	}
	// 更新user表的conversationNums
	err = l.UserModel.UpdateJianConversationNumByConversationIdAndUserId(l.ctx, userid)
	if err != nil {
		l.Logger.Errorf("%s database error.UpdateConversationNumByConversationIdAndUserId.info:%v", conversationModule, req)
		return biz.DBError
	}
	// 删除conversation下的所有message
	err = l.MessageModel.DeleteByConversationId(l.ctx, req.ConversationId)
	if err != nil {
		l.Logger.Errorf("%s database error.DeleteByConversationId.info:%v", conversationModule, req)
		return biz.DBError
	}
	return nil
}
