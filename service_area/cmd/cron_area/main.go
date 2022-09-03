package main

import (
	"context"

	"github.com/avast/retry-go/v4"
	"github.com/syukri21/mercari/common/helper"
	"github.com/syukri21/mercari/common/initialize"
	"github.com/syukri21/mercari/common/telemetry"
	"github.com/syukri21/mercari/service_area/cmd/config"
	"github.com/syukri21/mercari/service_area/repository/agent"
	"github.com/syukri21/mercari/service_area/repository/redis"
	"github.com/syukri21/mercari/service_area/usecase"

	"log"
	"os"
	"time"
)

var (
	filename      = "cmd/grpc/main.go"
	method        = "main"
	traceFilename = "cron_service_area_traces.txt"
)

func main() {
	ctx := context.Background()
	// Open Config
	serviceConfig, errConfig := config.InitConfig()
	helper.CheckError(errConfig)

	var cronRun = serviceConfig.Cron.FillArea
	if cronRun == "" {
		log.Fatal("crontab schedule is empty")
	}

	//	Open Telemetry
	l := log.New(os.Stdout, "", 0)
	f, err := os.Create(traceFilename)
	if err != nil {
		helper.CheckError(err)
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	tel := telemetry.NewTelemetry(f, l)
	tel.InitProvider()
	tel.StartTracerProvider()
	defer tel.Shutdown()

	//	Initialization
	redisPool := initialize.NewRedis(serviceConfig.Redis)

	//	Repository
	redisRepo := redis.NewRepositoryRedis(redisPool, tel.Log)
	agentRepo := agent.NewIKSRepository(tel.Log)

	uc := usecase.NewCronAreaUsecase(agentRepo, tel.Log, redisRepo)
	//c := cron.New()
	//
	//_, err = c.AddFunc(cronRun, func() {
	err = retry.Do(func() error {
		tel.Log.Printf("Running uc.CronFillArea ...")
		ctx = context.WithValue(ctx, "config", serviceConfig)
		err := uc.CronFillArea(ctx)
		return err
	}, retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
		// apply a default exponential back off strategy
		tel.Log.Printf("Retry uc.CronFillArea ...")

		return retry.BackOffDelay(n, err, config)
	}))
	//})

	tel.Log.Printf("Done uc.CronFillArea ...")

}
