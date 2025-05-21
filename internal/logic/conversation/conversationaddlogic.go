package conversation

import (
	"context"
	"encoding/json"
	"go-chat-sse/internal/biz"
	"go-chat-sse/internal/model"
	"go-chat-sse/internal/tools"

	"go-chat-sse/internal/svc"
	"go-chat-sse/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const conversationModule = "[conversationModule]"

type ConversationaddLogic struct {
	logx.Logger
	ctx                   context.Context
	svcCtx                *svc.ServiceContext
	ConversationModel     model.ConversationModel
	UserConversationModel model.UserConversationModel
	UserModel             model.UserModel
}

func NewConversationaddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConversationaddLogic {
	return &ConversationaddLogic{
		Logger:                logx.WithContext(ctx),
		ctx:                   ctx,
		svcCtx:                svcCtx,
		ConversationModel:     model.NewConversationModel(svcCtx.Mysql),
		UserConversationModel: model.NewUserConversationModel(svcCtx.Mysql),
		UserModel:             model.NewUserModel(svcCtx.Mysql),
	}
}

// Conversationadd 新增大会话
func (l *ConversationaddLogic) Conversationadd() (*types.ConversationAddResp, error) {

	// 获取userid
	userid, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		l.Logger.Infof("%s token error", conversationModule)
		return nil, biz.TokenError
	}
	// 生成conversationid，插入conversation相关表，并返回
	insertOne, err := l.ConversationModel.Insert(l.ctx, &model.Conversation{
		UserId:      userid,
		MessageNums: 0,
		FirstTime:   tools.GetNowTime(),
		LastTime:    tools.GetNowTime(),
	})
	if err != nil {
		l.Logger.Errorf("%s database error.info:%v", conversationModule, insertOne)
		return nil, biz.DBError
	}
	conversationId, err := insertOne.LastInsertId()
	if err != nil {
		l.Logger.Errorf("%s database error.info:%v", conversationModule, insertOne)
		return nil, biz.DBError
	}
	// 插入user_conversation表
	_, err = l.UserConversationModel.Insert(l.ctx, &model.UserConversation{
		ConversationId: conversationId,
		UserId:         userid,
		FirstTime:      tools.GetNowTime(),
		LastTime:       tools.GetNowTime(),
	})
	if err != nil {
		l.Logger.Errorf("%s database error.info:%v", conversationModule, insertOne)
		return nil, biz.DBError
	}
	// 更新用户表的conversation_nums
	err = l.UserModel.UpdateConversationNumByConversationIdAndUserId(l.ctx, userid)
	if err != nil {
		l.Logger.Errorf("%s database error.UpdateConversationNumByConversationIdAndUserId.info:%v", conversationModule, insertOne)
		return nil, biz.DBError
	}
	resp := &types.ConversationAddResp{
		ConversationId: conversationId,
	}

	l.Logger.Errorf("%s conversationadd success.info:%v", conversationModule, insertOne)
	logx.WithContext(l.ctx).Errorf("%s conversationadd success.info:%v", conversationModule, insertOne)

	return resp, nil
}
