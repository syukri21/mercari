package grpc

import (
	"context"

	grpctransport "github.com/go-kit/kit/transport/grpc"
)

// GrpcServer ...
type GrpcServer struct {
	register           grpctransport.Handler
	login              grpctransport.Handler
	refreshAccessToken grpctransport.Handler
	logout             grpctransport.Handler
	verifyRegister     grpctransport.Handler
	getLoginHistories  grpctransport.Handler
}

func (g GrpcServer) Register(ctx context.Context, request *RegisterRequest) (*RegisterResponse, error) {
	_, resp, err := g.register.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*RegisterResponse), nil
}

func (g GrpcServer) Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error) {
	_, resp, err := g.login.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*LoginResponse), nil
}

func (g GrpcServer) RefreshAccessToken(ctx context.Context, request *RefreshAccessTokenRequest) (*RefreshAccessTokenResponse, error) {
	_, resp, err := g.refreshAccessToken.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*RefreshAccessTokenResponse), nil
}

func (g GrpcServer) Logout(ctx context.Context, request *LogoutRequest) (*LogoutResponse, error) {
	_, resp, err := g.logout.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*LogoutResponse), nil
}

func (g GrpcServer) VerifyRegister(ctx context.Context, request *VerifyRegisterRequest) (*VerifyRegisterResponse, error) {
	_, resp, err := g.verifyRegister.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*VerifyRegisterResponse), nil
}

func (g GrpcServer) GetLoginHistories(ctx context.Context, request *GetLoginHistoriesRequest) (*GetLoginHistoriesResponse, error) {
	_, resp, err := g.getLoginHistories.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*GetLoginHistoriesResponse), nil
}

func (g GrpcServer) mustEmbedUnimplementedServiceAuthServer() {
	//TODO implement me
	panic("implement me")
}
