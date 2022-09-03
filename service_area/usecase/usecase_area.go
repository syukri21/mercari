package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/syukri21/mercari/common/helper"
	"github.com/syukri21/mercari/service_area/constant"
	"github.com/syukri21/mercari/service_area/model"
	"github.com/syukri21/mercari/service_area/repository"
	constantAuth "github.com/syukri21/mercari/service_auth/constant"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"log"
	"time"
)

// Area interface
type Area interface {
	GetAreaInfo(ctx context.Context, areaType string, key string) (model.AreaData, error)
}

// UsecaseArea  ...
type UsecaseArea struct {
	Logger *log.Logger
	Redis  repository.RedisRepository
	Agent  repository.Agent
}

func (a *UsecaseArea) GetAreaInfo(ctx context.Context, areaType string, key string) (result model.AreaData, err error) {
	_, span := otel.Tracer(constant.ServicesName).Start(ctx, "GetAreaInfo")
	start := time.Now()
	span.SetAttributes(attribute.String("Start Time", time.Now().Format(time.RFC3339)))

	defer func() {
		span.SetAttributes(attribute.String("End Time", time.Now().Format(time.RFC3339)))
		span.SetAttributes(attribute.String("Takes Time", time.Now().Sub(start).String()))
		span.End()
	}()

	result = model.AreaData{
		Key:   key,
		Value: "",
		Data:  make([]model.AreaRedis, 0),
	}

	infoAll, _, err := a.Redis.GetAreaInfo(ctx, areaType, key)
	if err != nil {
		span.RecordError(fmt.Errorf("a.Redis.GetAreaInfo err %v", err.Error()))
		span.SetStatus(codes.Error, err.Error())
		return model.AreaData{}, errors.New(constantAuth.StatusNotFound)
	}

	if infoAll == "null" {
		data, err := a.Agent.GetAreaData(ctx, key, areaType)
		if err != nil {
			span.RecordError(fmt.Errorf("a.Agent.GetAreaData err %v", err.Error()))
			span.SetStatus(codes.Error, err.Error())
			return model.AreaData{}, errors.New(constantAuth.StatusNotFound)
		}
		result.Data = data

		jsonData, _ := json.Marshal(result.Data)
		err = a.Redis.SaveAreaData(ctx, model.AreaRedis{
			Key:   helper.BuildKey(areaType, key),
			Value: string(jsonData),
		})
		if err != nil {
			span.RecordError(fmt.Errorf("a.Agent.GetAreaData err %v", err.Error()))
			span.SetStatus(codes.Error, err.Error())
			return model.AreaData{}, errors.New(constantAuth.StatusNotFound)
		}
	} else {
		err = json.Unmarshal([]byte(infoAll), &result.Data)
		if err != nil {
			return model.AreaData{}, errors.New(constantAuth.StatusNotFound)
		}
	}

	return result, nil

}

func NewAreaUsecase(logger *log.Logger, redis repository.RedisRepository, agent repository.Agent) Area {
	return &UsecaseArea{Logger: logger, Redis: redis, Agent: agent}
}
