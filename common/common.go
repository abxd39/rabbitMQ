package common

import (
	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
	Log "github.com/sirupsen/logrus"
	"sctek.com/typhoon/th-platform-gateway/common/worker"
)

var DB *xorm.Engine
var Config *ServerConfig
var RedisPool *redis.Client

//var Log *Logger

func LoadConfig() error {
	Config = &ServerConfig{}
	return Config.load()
}

var Pool *worker.Pool

func InitPool() {
	Pool = worker.NewPool(Config.MaxQueueSize)
	Pool.Run(Config.MaxWork)
	Log.Infof("goroutine的个数为%v,最大任务数为%v", Config.MaxWork, Config.MaxQueueSize)
}

func ClosePool() {
	Pool.Shutdown()
}

//func SetupLogger() error {
//	var err error
//	fmt.Printf("日志文件的路径为【%v】\r\n",Config.Log.LogFile)
//	Log, err = NewLogger(Config.Log.LogFile, Config.Log.TraceLevel)
//	return err
//}
