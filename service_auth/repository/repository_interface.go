//go:generate mockgen -source=repository_interface.go -destination=mock_repository.go -package=repository
package repository

import (
	"context"
	"github.com/syukri21/mercari/service_auth/model"
)

// JWTRepository ...
type JWTRepository interface {
	DecodeJWTToken(token string) (*model.CustomClaimData, error)
	GenerateJWTToken(plainToken model.PlainToken, jwtGet model.JWKGet) (string, error)
}

type PostgreRepository interface {
	CreateUser(ctx context.Context, request model.CreateUserRequest) (err error)
	ValidateUser(ctx context.Context, email string, pin string) (err error)
	GetUserByEmail(ctx context.Context, req model.LoginRequest) (user model.User, err error)
	CreateLoginHistory(ctx context.Context, req model.LoginHistory) error
	GetLoginHistories(ctx context.Context, req model.LoginHistoryRequest) ([]model.LoginHistory, error)
}

type RedisRepository interface {
	SaveLoginToken(ctx context.Context, refreshToken string, activeToken string, username string, email string, deviceId string) error
	GetLoginToken(ctx context.Context, email string) (model.RedisToken, error)
	ClearLoginToken(ctx context.Context, email string) error
}
