package setup

import (
	"github.com/go-redis/redis"
	remoteCfg "miaosha/pkg/config"
)

func InitRedis() error {
	client := redis.NewClient(&redis.Options{
		Addr:     remoteCfg.Redis.Host,
		Password: remoteCfg.Redis.Password,
		DB:       remoteCfg.Redis.Db,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return err
	}
	remoteCfg.Redis.RedisConn = client
	return nil
}
