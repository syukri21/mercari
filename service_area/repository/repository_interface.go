package repository

import (
	"context"
	"github.com/syukri21/mercari/service_area/model"
)

type RedisRepository interface {
	SaveAreaData(ctx context.Context, area model.AreaRedis) error
	GetAreaInfo(ctx context.Context, areaType string, key string, isAll bool) (name string, id string, err error)
}

type Agent interface {
	GetALLAreaData(ctx context.Context, saveFunc model.SaveAreaDataRedis) error
}
