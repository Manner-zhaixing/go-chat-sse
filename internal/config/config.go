package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	MysqlConfig MysqlConfig
	RedisConfig redis.RedisConf
	Auth        AuthConfig
	Snowflake   struct {
		WorkerId     int64
		DataCenterId int64
	}
	DeepSeek struct {
		ApiURL   string
		ApiKey   string
		ApiModel string
	}
	CommonUserIdStart int64
}

type MysqlConfig struct {
	DataSource     string
	ConnectTimeout int64
}

type AuthConfig struct {
	AccessSecret string
	Expire       int64
}

//type RedisConfig struct {
//	Host        string
//	Type        string
//	Pass        string
//	Tls         bool
//	NonBlock    bool
//	PingTimeout time.Duration
//}
