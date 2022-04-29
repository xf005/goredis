package cache

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/xf005/logger"
	"github.com/xf005/user/conf"
)

var (
	once   sync.Once
	client *redis.Client
)

func NewClient() *redis.Client {
	once.Do(func() {
		conf := conf.NewConf()
		client = redis.NewClient(&redis.Options{
			Addr:        conf.Redis.Host,
			Password:    conf.Redis.Pass,
			DB:          conf.Redis.Db,
			PoolSize:    5,                // 连接池大小
			MaxRetries:  3,                // 最大重试次数
			IdleTimeout: 10 * time.Second, // 空闲链接超时时间
			ReadTimeout: 5 * time.Second,  // 读取超时时间
		})
		pong, err := client.Ping(context.Background()).Result()
		if err != nil {
			logger.Error("pong err:", err.Error())
		}
		logger.Info(pong)
	})
	return client
}
