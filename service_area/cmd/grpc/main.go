package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/syukri21/mercari/common/helper"
	"github.com/syukri21/mercari/common/initialize"
	"github.com/syukri21/mercari/common/telemetry"
	"github.com/syukri21/mercari/service_area/cmd/config"
	"github.com/syukri21/mercari/service_area/repository/agent"
	"github.com/syukri21/mercari/service_area/repository/redis"
	areaGrpc "github.com/syukri21/mercari/service_area/transport/grpc"
	"github.com/syukri21/mercari/service_area/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	filename      = "cmd/grpc/main.go"
	method        = "main"
	traceFilename = "service_area_traces.txt"
)

func main() {
	ctx := context.Background()
	// Open Config
	serviceConfig, errConfig := config.InitConfig()
	helper.CheckError(errConfig)

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

	uc := usecase.NewAreaUsecase(tel.Log, redisRepo, agentRepo)

	// Initialize GRPC
	gRPCAddr := flag.String("grpc", fmt.Sprintf(":%d", serviceConfig.App.GRPCPort), "gRPC listen address")
	errChan := make(chan error)

	go func() {
		tel.Log.Printf("transport grpc")
		tel.Log.Printf("address %s", *gRPCAddr)
		tel.Log.Printf("gRPC server is listening")
		ctx = context.WithValue(ctx, "config", serviceConfig)

		listener, errListener := net.Listen("tcp", *gRPCAddr)
		if errListener != nil {
			errChan <- errListener
			return
		}

		handlerMaker := areaGrpc.NewHandlerMaker()
		handler := areaGrpc.MakeHandler(ctx, uc, handlerMaker)

		grpcServer := grpc.NewServer()
		areaGrpc.RegisterServiceAreaServer(grpcServer, handler)

		healthService := areaGrpc.NewHealthChecker()
		grpc_health_v1.RegisterHealthServer(grpcServer, healthService)

		errChan <- grpcServer.Serve(listener)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		errChan <- fmt.Errorf("%s", <-c)

		tel.Log.Printf("filename %s", filename)
		tel.Log.Printf("method %s", method)
		tel.Log.Printf("Gracefully Stop")
	}()

	log.Println(<-errChan)
}
