package redis

import (
	"context"
	"github.com/syukri21/mercari/common/helper"
	"github.com/syukri21/mercari/service_area/constant"
	"github.com/syukri21/mercari/service_area/model"
	"github.com/syukri21/mercari/service_area/repository"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// RepositoryRedis ...
type RepositoryRedis struct {
	Client *redis.Client
	l      *log.Logger
}

func (r RepositoryRedis) SaveAreaData(ctx context.Context, area model.AreaRedis) error {
	err := r.Client.HSet(ctx, constant.AreaRedisKey, area.Key, area.Value).Err()
	if err != nil {
		return err
	}
	err = r.Client.Expire(ctx, constant.AreaRedisKey, time.Hour*30*24).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r RepositoryRedis) GetAreaInfo(ctx context.Context, areaType string, key string) (name string, id string, err error) {
	key = helper.GetKey(areaType, key)
	hget := r.Client.HGet(ctx, constant.AreaRedisKey, key)

	var value string
	err = hget.Scan(&value)
	if err != nil {
		return "", "", err
	}

	return value, key, err
}

func NewRepositoryRedis(client *redis.Client, l *log.Logger) repository.RedisRepository {
	return &RepositoryRedis{Client: client, l: l}
}
