package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"go-chat-sse/internal/biz"
	"go-chat-sse/internal/model"

	"go-chat-sse/internal/svc"
	"go-chat-sse/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const userInfoModule = "[userInfo]"

type GetUserInfoLogic struct {
	logx.Logger
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	userModel model.UserModel
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger:    logx.WithContext(ctx),
		ctx:       ctx,
		svcCtx:    svcCtx,
		userModel: model.NewUserModel(svcCtx.Mysql),
	}
}

// GetUserInfo 获取用户信息
func (l *GetUserInfoLogic) GetUserInfo() (resp *types.UserInfoResp, err error) {
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		l.Logger.Info("%s token error", userInfoModule)
		return nil, biz.TokenError
	}

	user, err := l.userModel.FindOne(l.ctx, userId)
	if err != nil && errors.Is(err, model.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
		l.Logger.Infof("%s user not found.info:%v", biz.UserNotFound, userId)
		return nil, biz.UserNotFound
	}
	if err != nil {
		l.Logger.Errorf("%s database error.info:%v", userInfoModule, userId)
		return nil, biz.DBError
	}

	l.Logger.Info("%s get user info success.info:%v", userInfoModule, userId)
	return &types.UserInfoResp{
		Id:       user.Id,
		Username: user.Username,
	}, nil
}
