package user

import (
	"context"
	"go-chat-sse/internal/biz"
	"go-chat-sse/internal/model"
	"go-chat-sse/internal/tools"

	"go-chat-sse/internal/svc"
	"go-chat-sse/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const userRegisterModule = "[register]"

type RegisterLogic struct {
	logx.Logger
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	UserModel model.UserModel
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger:    logx.WithContext(ctx),
		ctx:       ctx,
		svcCtx:    svcCtx,
		UserModel: model.NewUserModel(svcCtx.Mysql),
	}
}

// checkData 校验数据
func (l *RegisterLogic) checkData(req *types.RegisterReq) error {
	if req.Username == "" || req.Password == "" || len(req.Username) > 20 || len(req.Password) > 20 {
		return biz.ParamError
	}
	return nil
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (*types.RegisterResp, error) {
	// 1.校验username和password数据是否
	err := l.checkData(req)
	if err != nil {
		l.Logger.Infof("%s username or password error.info:%v", userRegisterModule, *req)
		return nil, biz.ParamError
	}
	// 2.根据username查询用户表，判断password是否对劲
	userOne, err := l.UserModel.FindByUsername(l.ctx, req.Username)
	if err != nil {
		l.Logger.Errorf("%s database error.info:%v", userRegisterModule, *req)
		return nil, biz.DBError
	}
	if userOne != nil {
		// 注册过了，不允许重复注册
		l.Logger.Infof("%s user already exists.info:%v", userRegisterModule, *req)
		return nil, biz.AlreadyRegister
	}
	// 3.如果不存在，用户表插入用户数据，并返回成功(密码要加密)
	passwordCry, err := tools.GetHashPassword(req.Password)
	if err != nil {
		l.Logger.Infof("%s password error,加密失败.info:%v", userRegisterModule, *req)
	}
	_, err = l.UserModel.Insert(l.ctx, &model.User{
		Username:         req.Username,
		Password:         passwordCry,
		ConversationNums: 0,
		RegisterTime:     tools.GetNowTime(),
		LastLoginTime:    tools.GetNowTime(),
	})
	if err != nil {
		l.Logger.Errorf("%s database error.info:%v", userRegisterModule, *req)
		return nil, biz.DBError
	}

	l.Logger.Infof("%s register success.info:%v", userRegisterModule, *req)
	return &types.RegisterResp{
		Msg: "注册成功",
	}, nil
}
