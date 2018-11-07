package common

import (
	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
)

var DB *xorm.Engine
var Config *ServerConfig
var RedisPool *redis.Client
var Log *Logger

func LoadConfig() error {
	Config = &ServerConfig{}
	return Config.load()
}

func SetupLogger() error {
	var err error
	Log, err = NewLogger(Config.Log.LogFile, Config.Log.TraceLevel)
	return err
}
