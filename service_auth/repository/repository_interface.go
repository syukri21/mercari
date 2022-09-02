//go:generate mockgen -source=repository_interface.go -destination=mock_repository.go -package=repository
package repository

import "github.com/syukri21/mercari/service_auth/model"

// JWTRepository ...
type JWTRepository interface {
	DecodeJWTToken(token string) (*model.CustomClaimData, error)
	GenerateJWTToken(plainToken model.PlainToken, jwtGet model.JWKGet) (string, error)
}
