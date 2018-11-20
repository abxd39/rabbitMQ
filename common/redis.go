package common

import (
	"time"

	"github.com/go-redis/redis"
)

const (
	RedisKeyOpenUser    string        = "open:u:%s"
	DefaultTimeDuration time.Duration = time.Minute * 30
)

func OpenRedis() error {
	RedisPool = redis.NewClient(&redis.Options{
		Addr:     Config.Redis.Address,
		Password: Config.Redis.PassWord,
		DB:       Config.Redis.Database,
	})

	_, err := RedisPool.Ping().Result()
	//fmt.Println("redis initial status:", pong, ";err:", err)
	return err
}

func SetRedisValue(key string, val interface{}) error {
	err := RedisPool.Set(key, val, time.Duration(DefaultTimeDuration)).Err()
	if err != nil {
		Log.Errorf("redis setvalue failed:", err)
	}
	return err
}

func SetRedisValueWithDuration(key string, val interface{}, duration time.Duration) error {
	err := RedisPool.Set(key, val, duration).Err()
	if err != nil {
		Log.Errorf("redis setvalue failed:", err)
	}
	return err
}

func GetRedisValue(key string) (string, error) {
	val, err := RedisPool.Get(key).Result()
	if err == redis.Nil {
		Log.Errorf("redis key ", key, " does not exist")
		return "", nil
	} else if err != nil {
		Log.Errorf("redis getvalue failed:", err)
	}
	return val, err
}
