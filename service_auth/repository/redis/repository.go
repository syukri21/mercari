package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/syukri21/mercari/service_auth/model"
	"github.com/syukri21/mercari/service_auth/repository"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// RepositoryRedis ...
type RepositoryRedis struct {
	Client *redis.Client
	l      *log.Logger
}

// NewRepositoryRedis ...
func NewRepositoryRedis(l *log.Logger, rds *redis.Client) repository.RedisRepository {
	return &RepositoryRedis{
		Client: rds,
		l:      l,
	}
}

func (r *RepositoryRedis) SaveLoginToken(ctx context.Context, refreshToken string, activeToken string, username string, email string, deviceId string) error {
	key := fmt.Sprintf("login_token_%s", email)
	cmds := map[string]interface{}{
		"device_id":     deviceId,
		"username":      username,
		"email":         email,
		"refresh_token": refreshToken,
		"active_token":  activeToken,
	}

	err := r.Client.HSet(ctx, key, cmds).Err()
	if err != nil {
		return errors.New(fmt.Sprintf("[SaveLoginToken] error when do HMSET to save login token. err : %s", err.Error()))
	}
	err = r.Client.Expire(ctx, key, 300*time.Second).Err()
	if err != nil {
		return errors.New(fmt.Sprintf("[SaveLoginToken] error when do EXPIRE to set expiration redis. err : %s", err.Error()))
	}

	return nil
}

func (r *RepositoryRedis) GetLoginToken(ctx context.Context, email string) (model.RedisToken, error) {
	result := model.RedisToken{}

	key := fmt.Sprintf("login_token_%s", email)
	redisMap, err := r.Client.HGetAll(ctx, key).Result()
	if len(redisMap) == 0 {
		return result, err
	}

	result.DeviceID = redisMap["device_id"]
	result.Email = redisMap["device_id"]
	result.Username = redisMap["username"]
	result.RefreshToken = redisMap["refresh_token"]
	result.ActiveToken = redisMap["active_token"]

	return result, nil
}

func (r *RepositoryRedis) ClearLoginToken(ctx context.Context, email string) error {
	key := fmt.Sprintf("login_token_%s", email)
	return r.Client.Del(ctx, key).Err()
}
