package usecase

import (
	"context"
	"github.com/syukri21/mercari/service_auth/model"
)

// AuthUsecase interface
type AuthUsecase interface {
	Register(ctx context.Context, req model.RegisterRequest) (model.RegisterResponse, error)
	Login(ctx context.Context, req model.LoginRequest) (model.LoginResponse, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (model.RefreshAccessToken, error)
	Logout(ctx context.Context, req model.LogoutRequest) (model.LogoutResponse, error)
}
