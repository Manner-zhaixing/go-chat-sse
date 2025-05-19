package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-chat-sse/internal/config"
	"go-chat-sse/internal/db"
	"go-chat-sse/internal/tools"
)

type ServiceContext struct {
	Config        config.Config
	Mysql         sqlx.SqlConn
	Redis         *redis.Redis
	IdWorker      *tools.Worker // ID生成器
	IdWorkerRedis *tools.SimpleRedisIDGenerator
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := db.NewMysql(c.MysqlConfig)
	redisConn := db.NewRedisConn(c.RedisConfig)
	IdWorker, err := tools.NewWorker(c.Snowflake.WorkerId)
	IdWorkerRedis, err := tools.NewSimpleRedisIDGenerator(redisConn, c.Snowflake.WorkerId)
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:        c,
		Mysql:         sqlConn,
		Redis:         redisConn,
		IdWorker:      IdWorker,
		IdWorkerRedis: IdWorkerRedis,
	}
}
