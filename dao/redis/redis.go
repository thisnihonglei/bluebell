package redis

import (
	"bluebell/setting"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
)

func Init(cfg *setting.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password, // 密码
		DB:       cfg.DB,       // 数据库
		PoolSize: cfg.PoolSize, // 连接池大小
	})

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return nil
	}
	return
}

func Close() {
	_ = client.Close()
}
