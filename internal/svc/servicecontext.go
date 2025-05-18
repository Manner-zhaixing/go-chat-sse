package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-chat-sse/internal/config"
	"go-chat-sse/internal/db"
)

type ServiceContext struct {
	Config config.Config
	Mysql  sqlx.SqlConn
	Redis  *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := db.NewMysql(c.MysqlConfig)
	redisConn := db.NewRedisConn(c.RedisConfig)
	return &ServiceContext{
		Config: c,
		Mysql:  sqlConn,
		Redis:  redisConn,
	}
}
