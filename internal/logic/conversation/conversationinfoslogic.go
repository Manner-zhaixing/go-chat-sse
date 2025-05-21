package conversation

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-chat-sse/internal/biz"
	"go-chat-sse/internal/model"

	"go-chat-sse/internal/svc"
	"go-chat-sse/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const conversationInfosModule = "[conversationInfosModule]"

type ConversationinfosLogic struct {
	logx.Logger
	ctx               context.Context
	svcCtx            *svc.ServiceContext
	ConversationModel model.ConversationModel
}

func NewConversationinfosLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConversationinfosLogic {
	return &ConversationinfosLogic{
		Logger:            logx.WithContext(ctx),
		ctx:               ctx,
		svcCtx:            svcCtx,
		ConversationModel: model.NewConversationModel(svcCtx.Mysql),
	}
}

func (l *ConversationinfosLogic) Conversationinfos() (*types.ConversationInfosResp, error) {
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		l.Logger.Infof("%s token error", conversationInfosModule)
		return nil, biz.TokenError
	}
	conversationInfos, err := l.ConversationModel.FindOneByUserId(l.ctx, userId)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			l.Logger.Infof("%s userId:%d not found", conversationInfosModule, userId)
			return nil, nil
		} else {
			l.Logger.Errorf("%s db error,err:%s", conversationInfosModule)
			return nil, biz.DBError
		}
	}
	var temp []types.ConversationInfo
	for _, conversationInfo := range *conversationInfos {
		temp = append(temp, types.ConversationInfo{
			ConversationId: conversationInfo.Id,
			FirstTime:      conversationInfo.FirstTime.Format("2006-01-02 15:04:05"),
			LastTime:       conversationInfo.LastTime.Format("2006-01-02 15:04:05"),
			UserId:         conversationInfo.UserId,
		})
	}
	return &types.ConversationInfosResp{
		ConversationInfos: temp,
	}, nil
}
