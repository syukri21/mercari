package usecase

import (
	"context"
	"github.com/syukri21/mercari/service_area/constant"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"log"
	"time"

	"github.com/syukri21/mercari/service_area/repository"
)

// CronArea  ...
type CronArea interface {
	CronFillArea(ctx context.Context) error
}

type CronAreaUsecase struct {
	AreaAgent repository.Agent
	Redis     repository.RedisRepository
	Logger    *log.Logger
}

func (c CronAreaUsecase) CronFillArea(ctx context.Context) error {
	_, span := otel.Tracer(constant.ServicesName).Start(ctx, "CronFillArea")
	defer span.End()
	start := time.Now()

	c.Logger.Printf("[CronFillArea Start] %v", start.Format(time.RFC3339))
	span.SetAttributes(attribute.String("Start Time", start.Format(time.RFC3339)))

	defer func() {
		c.Logger.Printf("[CronFillArea End] %v takes time", start.Format(time.RFC3339), time.Now().Sub(start))
		span.SetAttributes(attribute.String("End Time", time.Now().Format(time.RFC3339)))
		span.SetAttributes(attribute.String("Takes Time", time.Now().Sub(start).String()))
	}()

	err := c.AreaAgent.GetALLAreaData(ctx, c.Redis.SaveAreaData)
	if err != nil {
		c.Logger.Printf("[Error when c.AreaAgent.GetALLAreaData; %s]", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	c.Logger.Printf("[CronFillArea Success]", err)
	span.SetStatus(codes.Ok, "Success")
	return nil
}

func NewCronAreaUsecase(
	AreaAgent repository.Agent,
	Logger *log.Logger,
	Redis repository.RedisRepository,
) CronArea {
	return &CronAreaUsecase{AreaAgent: AreaAgent, Logger: Logger, Redis: Redis}
}
