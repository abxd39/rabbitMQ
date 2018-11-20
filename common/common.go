package common

import (
	"fmt"
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
	fmt.Printf("日志文件的路径为【%v】\r\n",Config.Log.LogFile)
	Log, err = NewLogger(Config.Log.LogFile, Config.Log.TraceLevel)
	return err
}
