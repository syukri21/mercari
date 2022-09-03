package grpc

import (
	"context"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/syukri21/mercari/service_auth/constant"
	"github.com/syukri21/mercari/service_auth/helper"
	"github.com/syukri21/mercari/service_auth/model"
	"github.com/syukri21/mercari/service_auth/transport/validation"
	"github.com/syukri21/mercari/service_auth/usecase"
	"google.golang.org/grpc/health/grpc_health_v1"
	"time"
)

// HandlerMaker ...
type HandlerMaker struct {
}

func NewHandlerMaker() *HandlerMaker {
	return &HandlerMaker{}
}

// MakeHandler ...
func MakeHandler(oldCtx context.Context, usecase usecase.Auth, handlerMaker *HandlerMaker) ServiceAuthServer {
	validator := validation.NewValidator()
	return &GrpcServer{
		register: kitgrpc.NewServer(
			func(ctx context.Context, request interface{}) (response interface{}, err error) {
				ctx = context.WithValue(ctx, "config", oldCtx.Value("config").(model.Config))
				req := request.(model.RegisterRequest)
				err = validator.Struct(req)
				if err != nil {
					return nil, err
				}
				return usecase.Register(ctx, req)
			},
			handlerMaker.decodeRegisterRequest,
			handlerMaker.encodeRegisterResponse,
		),
		login: kitgrpc.NewServer(
			func(ctx context.Context, request interface{}) (response interface{}, err error) {
				ctx = context.WithValue(ctx, "config", oldCtx.Value("config").(model.Config))
				req := request.(model.LoginRequest)
				err = validator.Struct(req)
				if err != nil {
					return nil, err
				}
				return usecase.Login(ctx, req)
			},
			handlerMaker.decodeLoginRequest,
			handlerMaker.encodeLoginResponse,
		),
		refreshAccessToken: kitgrpc.NewServer(
			func(ctx context.Context, request interface{}) (response interface{}, err error) {
				ctx = context.WithValue(ctx, "config", oldCtx.Value("config").(model.Config))
				req := request.(*RefreshAccessTokenRequest)
				session, err := helper.CheckSession(ctx)
				if err != nil {
					return nil, err
				}
				ctx = context.WithValue(ctx, "session", session)
				accessToken, err := usecase.RefreshAccessToken(ctx, req.DeviceId, req.Email, req.ActiveAccessToken)
				return accessToken, err
			},
			handlerMaker.decodeRefreshAccessTokenRequest,
			handlerMaker.encodeRefreshAccessTokenResponse,
		),
		logout: kitgrpc.NewServer(
			func(ctx context.Context, request interface{}) (response interface{}, err error) {
				ctx = context.WithValue(ctx, "config", oldCtx.Value("config").(model.Config))
				req := request.(*LogoutRequest)
				err = usecase.Logout(ctx, req.Email, req.ActiveAccessToken)
				return nil, err
			},
			handlerMaker.decodeLogoutRequest,
			handlerMaker.encodeLogoutResponse,
		),
		verifyRegister: kitgrpc.NewServer(
			func(ctx context.Context, request interface{}) (response interface{}, err error) {
				ctx = context.WithValue(ctx, "config", oldCtx.Value("config").(model.Config))
				req := request.(model.VerifyRegisterRequest)
				err = validator.Struct(req)
				if err != nil {
					return nil, err
				}
				register, err := usecase.VerifyRegister(ctx, req)
				return register, err
			},
			handlerMaker.decodeVerifyRegisterRequest,
			handlerMaker.encodeVerifyRegisterResponse,
		),
		getLoginHistories: kitgrpc.NewServer(
			func(ctx context.Context, request interface{}) (response interface{}, err error) {
				ctx = context.WithValue(ctx, "config", oldCtx.Value("config").(model.Config))
				session, err := helper.CheckSession(ctx)
				if err != nil {
					return nil, err
				}
				req := request.(model.LoginHistoryRequest)
				ctx = context.WithValue(ctx, "session", session)
				req.Email = session.Email
				return usecase.GetLoginHistories(ctx, req)
			},
			handlerMaker.decodeGetLoginHistories,
			handlerMaker.encodeGetLoginHistories,
		),
	}
}

// decodeRegister ...
func (h *HandlerMaker) decodeRegisterRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*RegisterRequest)
	request := model.RegisterRequest{
		Email:           req.Email,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	}
	return request, nil
}

// encodeRegister ...
func (h *HandlerMaker) encodeRegisterResponse(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(model.RegisterResponse)
	return &RegisterResponse{
		Email:  req.Email,
		Status: constant.StatusOK,
		Error:  "",
	}, nil
}

// decodeLogin ...
func (h *HandlerMaker) decodeLoginRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*LoginRequest)
	request := model.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
		DeviceID: req.DeviceId,
	}
	return request, nil
}

// encodeLogin ...
func (h *HandlerMaker) encodeLoginResponse(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(model.LoginResponse)
	return &LoginResponse{
		Email:             req.Email,
		Username:          req.Username,
		ActiveAccessToken: req.ActiveAccessToken,
		Status:            constant.StatusOK,
		Error:             "",
	}, nil
}

// decodeRefreshAccessTokenRequest ...
func (h *HandlerMaker) decodeRefreshAccessTokenRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*RefreshAccessTokenRequest)
	return req, nil
}

// encodeRefreshAccessTokenResponse ...
func (h *HandlerMaker) encodeRefreshAccessTokenResponse(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(model.RefreshAccessToken)
	return &RefreshAccessTokenResponse{
		ActiveAccessToken: req.AccessToken,
		Status:            constant.StatusOK,
		Error:             "",
	}, nil
}

// decodeLogoutRequest ...
func (h *HandlerMaker) decodeLogoutRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*LogoutRequest)
	return req, nil
}

// encodeLogoutResponse ...
func (h *HandlerMaker) encodeLogoutResponse(ctx context.Context, r interface{}) (interface{}, error) {
	return &LogoutResponse{
		Status: constant.StatusOK,
		Error:  "",
	}, nil
}

// decodeVerifyRegisterRequest ...
func (h *HandlerMaker) decodeVerifyRegisterRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*VerifyRegisterRequest)
	return model.VerifyRegisterRequest{
		Email:       req.Email,
		ActivateKey: req.ActivateKey,
	}, nil
}

// encodeVerifyRegisterResponse ...
func (h *HandlerMaker) encodeVerifyRegisterResponse(ctx context.Context, r interface{}) (interface{}, error) {
	return &VerifyRegisterResponse{
		Email:       "",
		ActivateKey: "",
	}, nil
}

// decodeVerifyRegisterRequest ...
func (h *HandlerMaker) decodeGetLoginHistories(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*GetLoginHistoriesRequest)
	var offset int
	if req.Page == 1 {
		offset = 0
	} else {
		offset = int(req.Page) * int(req.Limit)
	}
	return model.LoginHistoryRequest{
		Limit:  int(req.Limit),
		Offset: offset,
	}, nil
}

// encodeVerifyRegisterResponse ...
func (h *HandlerMaker) encodeGetLoginHistories(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(model.LoginHistoryResponse)
	var data []*GetLoginHistoriesResponse_Data
	for _, datum := range req.Data {
		data = append(data, &GetLoginHistoriesResponse_Data{
			Email:    datum.Email,
			Username: datum.Username,
			DeviceId: datum.DeviceId,
			LoginAt:  datum.LoginAt.Format(time.RFC3339),
		})
	}

	return &GetLoginHistoriesResponse{
		Data:   data,
		Status: constant.StatusOK,
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
