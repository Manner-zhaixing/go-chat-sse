package chat

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

const messageModule = "[messageModule]"

type MessageLogic struct {
	logx.Logger
	ctx          context.Context
	svcCtx       *svc.ServiceContext
	messageModel model.MessageModel
	SessionModel model.SessionModel
}

func NewMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageLogic {
	return &MessageLogic{
		Logger:       logx.WithContext(ctx),
		ctx:          ctx,
		svcCtx:       svcCtx,
		messageModel: model.NewMessageModel(svcCtx.Mysql),
		SessionModel: model.NewSessionModel(svcCtx.Mysql),
	}
}

func (l *MessageLogic) checkData(req *types.MessageReq) error {
	if req.ConversationId <= 0 || req.ModelId <= 0 || req.FromId <= 0 || req.ToId <= 0 || req.Content == "" {
		return biz.ParamError
	}
	return nil
}

func (l *MessageLogic) Message(req *types.MessageReq) (*types.MessageResp, error) {
	// 1.校验数据
	err := l.checkData(req)
	if err != nil {
		l.Logger.Infof("%s checkData error: %v", messageModule, err)
		return nil, biz.ParamError
	}
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		l.Logger.Info("%s token error", messageModule)
		return nil, biz.TokenError
	}
	// 2.在message表中存储消息
	// 生成messageid
	messageId := l.svcCtx.IdWorker.Next()
	_, err = l.messageModel.Insert(l.ctx, &model.Message{
		MessageId:      messageId,
		UserId:         userId,
		ConversationId: req.ConversationId,
		ModelId:        int64(req.ModelId),
		FromId:         req.FromId,
		ToId:           req.ToId,
		Content:        req.Content,
		CurTime:        tools.GetNowTime(),
	})
	if err != nil {
		// 插入失败
		l.Logger.Errorf("%s database error.info:%v", messageModule, *req)
		return nil, biz.DBError
	}
	// 3.插入session表，获取session_id
	sessionid, _ := l.svcCtx.IdWorkerRedis.GenerateID()
	_, err = l.SessionModel.Insert(l.ctx, &model.Session{
		SessionId:      sessionid,
		ConversationId: req.ConversationId,
		UserId:         userId,
		MessageId:      messageId,
	})
	if err != nil {
		// 插入失败
		l.Logger.Errorf("%s database error.info:%v", messageModule, *req)
		return nil, biz.DBError
	}
	// 4.返回sessionid
	l.Logger.Infof("%s message success.info:%v", messageModule, *req)
	return &types.MessageResp{
		SessionId: sessionid,
	}, nil
}
