package db

import "github.com/zeromicro/go-zero/core/stores/redis"

func NewRedisConn(con redis.RedisConf) *redis.Redis {
	return redis.MustNewRedis(con)
}
