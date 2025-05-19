package user

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-chat-sse/internal/biz"
	"go-chat-sse/internal/model"
	"go-chat-sse/internal/tools"

	"go-chat-sse/internal/svc"
	"go-chat-sse/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const userLoginModule = "[login]"

type LoginLogic struct {
	logx.Logger
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	userModel model.UserModel
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger:    logx.WithContext(ctx),
		ctx:       ctx,
		svcCtx:    svcCtx,
		userModel: model.NewUserModel(svcCtx.Mysql),
	}
}

func (l *LoginLogic) checkData(req *types.LoginReq) error {
	if req.Username == "" || req.Password == "" || len(req.Username) > 20 || len(req.Password) > 20 {
		return biz.ParamError
	}
	return nil
}

func (l *LoginLogic) Login(req *types.LoginReq) (*types.LoginResp, error) {
	// 1.校验请求数据
	err := l.checkData(req)
	if err != nil {
		l.Logger.Infof("%s username or password error.info:%v", userLoginModule, *req)
		return nil, biz.ParamError
	}
	// 2.根据用户名查询userOne
	userOne, err := l.userModel.FindByUsername(l.ctx, req.Username)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			l.Logger.Infof("%s user not exists.info:%v", userLoginModule, *req)
			return nil, biz.UserNotFound
		} else {
			l.Logger.Errorf("%s database error.info:%v", biz.DBError, *req)
			return nil, biz.DBError
		}
	}
	// 3.校验密码是否正确
	if !tools.CheckPasswordHash(req.Password, userOne.Password) {
		// 密码错误
		l.Logger.Infof("%s password error.info:%v", userLoginModule, *req)
		return nil, biz.PasswordError
	}
	// 4.生成token返回
	token, err := biz.GetJwtToken(l.svcCtx.Config.Auth.AccessSecret, tools.GetNowTime().Unix(), l.svcCtx.Config.Auth.Expire, userOne.Id)
	if err != nil {
		l.Logger.Infof("%s token error.info:%v", userLoginModule, *req)
		return nil, biz.TokenGenError
	}

	l.Logger.Infof("%s login success.info:%v", userLoginModule, *req)
	return &types.LoginResp{
		Token: token,
	}, nil
}
