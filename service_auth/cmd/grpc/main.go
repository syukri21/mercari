package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/syukri21/mercari/common/helper"
	"github.com/syukri21/mercari/common/initialize"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/syukri21/mercari/common/telemetry"
	"github.com/syukri21/mercari/service_auth/cmd/config"
	"github.com/syukri21/mercari/service_auth/repository/jwt"
	"github.com/syukri21/mercari/service_auth/repository/postgre"
	"github.com/syukri21/mercari/service_auth/repository/redis"
	authGrpc "github.com/syukri21/mercari/service_auth/transport/grpc"
	"github.com/syukri21/mercari/service_auth/usecase"
	"google.golang.org/grpc"
)

var (
	filename      = "cmd/grpc/main.go"
	method        = "main"
	traceFilename = "service_auth_traces.txt"
)

func main() {
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
	dbPostgrePool := initialize.NewPostSqlServer(serviceConfig.Postgre)

	//	Repository
	postgreRepo := postgre.NewPostgreRepository(tel.Log, dbPostgrePool)
	redisRepo := redis.NewRepositoryRedis(tel.Log, redisPool)
	jwtRepo := jwt.NewJWTRepository()

	usecaseAuth := usecase.NewAuthUsecase(postgreRepo, jwtRepo, redisRepo, tel.Log, serviceConfig)

	// Initialize GRPC
	gRPCAddr := flag.String("grpc", fmt.Sprintf(":%d", serviceConfig.App.GRPCPort), "gRPC listen address")
	errChan := make(chan error)

	go func() {
		tel.Log.Printf("transport grpc")
		tel.Log.Printf("address %s", *gRPCAddr)
		tel.Log.Printf("gRPC server is listening")
		ctx := context.Background()
		ctx = context.WithValue(ctx, "config", serviceConfig)

		listener, errListener := net.Listen("tcp", *gRPCAddr)
		if errListener != nil {
			errChan <- errListener
			return
		}

		handlerMaker := authGrpc.NewHandlerMaker()
		handler := authGrpc.MakeHandler(ctx, usecaseAuth, handlerMaker)

		grpcServer := grpc.NewServer()
		authGrpc.RegisterServiceAuthServer(grpcServer, handler)

		healthService := authGrpc.NewHealthChecker()
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
