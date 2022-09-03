package initialize

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/syukri21/mercari/common/model"
)

// NewRedis ...
func NewRedis(redisConfig model.RedisConfig) *redis.Client {
	opt := &redis.Options{
		Addr:         fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		IdleTimeout:  time.Duration(redisConfig.IdleTimeout) * time.Second,
		MinIdleConns: redisConfig.MaxIdle / 2,
		PoolSize:     redisConfig.MaxActive,
	}
	c := redis.NewClient(opt)

	if err := c.Ping(context.TODO()).Err(); err != nil {
		panic(err)
	}

	return c
}
