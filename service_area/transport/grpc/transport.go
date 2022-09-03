package grpc

import (
	"context"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/syukri21/mercari/service_area/model"
	"github.com/syukri21/mercari/service_area/usecase"
	"github.com/syukri21/mercari/service_auth/transport/validation"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// HandlerMaker ...
type HandlerMaker struct {
}

func NewHandlerMaker() *HandlerMaker {
	return &HandlerMaker{}
}

// MakeHandler ...
func MakeHandler(oldCtx context.Context, usecase usecase.Area, handlerMaker *HandlerMaker) ServiceAreaServer {
	validator := validation.NewValidator()
	return &GrpcServer{
		getAreaInfo: kitgrpc.NewServer(
			func(ctx context.Context, request interface{}) (response interface{}, err error) {
				ctx = context.WithValue(ctx, "config", oldCtx.Value("config").(model.Config))
				req := request.(model.GetAreaInfoRequest)
				err = validator.Struct(req)
				if err != nil {
					return nil, err
				}
				return usecase.GetAreaInfo(ctx, req.AreaType, req.Key)
			},
			handlerMaker.decodeGetAreaInfoRequest,
			handlerMaker.encodeGetAreaInfoResponse,
		),
	}
}

// decodeGetAreaInfoRequest ...
func (h *HandlerMaker) decodeGetAreaInfoRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*GetAreaRequest)
	request := model.GetAreaInfoRequest{
		AreaType: req.AreaType,
		Key:      req.Key,
	}
	return request, nil
}

// encodeGetAreaInfoResponse ...
func (h *HandlerMaker) encodeGetAreaInfoResponse(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(model.AreaData)
	data := make([]*GetAreaInfoResponse_Data, 0)

	for _, value := range req.Data {
		data = append(data, &GetAreaInfoResponse_Data{
			Name: value.Value,
			Id:   value.Key,
		})
	}
	return &GetAreaInfoResponse{
		Name: req.Value,
		Id:   req.Key,
		Data: data,
	}, nil
}

// healthchecker ...
type healthchecker struct {
}

func NewHealthChecker() *healthchecker {
	return &healthchecker{}
}

func (h healthchecker) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (h healthchecker) Watch(req *grpc_health_v1.HealthCheckRequest, server grpc_health_v1.Health_WatchServer) error {
	return server.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
}
