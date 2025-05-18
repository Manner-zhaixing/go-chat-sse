package biz

const Ok = 200

var (
	DBError         = NewError(10000, "数据库错误")
	ParamError      = NewError(10001, "参数错误")
	AlreadyRegister = NewError(10100, "用户已注册")
	PasswordError   = NewError(10200, "密码错误")
	UserNotFound    = NewError(10300, "用户不存在")
	TokenGenError   = NewError(10400, "token错误")
	TokenError      = NewError(10500, "token过期")
	RedisErr        = NewError(10600, "redis错误")
)
