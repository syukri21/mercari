package usecase

import (
	"context"
	"log"

	"github.com/syukri21/mercari/service_auth/model"
	"github.com/syukri21/mercari/service_auth/repository"
)

// Auth interface
type Auth interface {
	Register(ctx context.Context, req model.RegisterRequest) (model.RegisterResponse, error)
	Login(ctx context.Context, req model.LoginRequest) (model.LoginResponse, error)
	RefreshAccessToken(ctx context.Context, deviceId string, email string, activeToken string) (result model.RefreshAccessToken, err error)
	Logout(ctx context.Context, email string, activeToken string) (err error)
	VerifyRegister(ctx context.Context, req model.VerifyRegisterRequest) (result model.VerifyRegisterResponse, err error)
	GetLoginHistories(ctx context.Context, req model.LoginHistoryRequest) (result model.LoginHistoryResponse, err error)
	CheckSession(ctx context.Context) (result model.CustomClaimData, err error)
}

type AuthUsecase struct {
	Postgre repository.PostgreRepository
	JWT     repository.JWTRepository
	Redis   repository.RedisRepository
	Logger  *log.Logger
	config  model.Config
}

func NewAuthUsecase(
	Postgre repository.PostgreRepository,
	JWT repository.JWTRepository,
	Redis repository.RedisRepository,
	Logger *log.Logger,
	config model.Config,
) Auth {
	return &AuthUsecase{Postgre: Postgre, JWT: JWT, Redis: Redis, Logger: Logger, config: config}
}
