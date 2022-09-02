package model

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type CustomClaim struct {
	Type             string          `json:"typ,omitempty"`
	RefreshExpiresAt int64           `json:"-"`
	IssuedTime       *time.Time      `json:"-"`
	Data             CustomClaimData `json:"data"`
	jwt.StandardClaims
}

type CustomClaimData struct {
	Username  string `json:"use,omitempty"`
	Email     string `json:"ema,omitempty"`
	DeviceID  string `json:"dvc,omitempty"`
	SessionID string `json:"sid,omitempty"`
}

// PlainToken ...
type PlainToken struct {
	Aud   string          `json:"aud,omitempty"`
	Data  CustomClaimData `json:"data,omitempty"`
	Did   string          `json:"did,omitempty"`
	Exp   int32           `json:"exp,omitempty"`
	Iss   string          `json:"iss,omitempty"`
	Jti   string          `json:"jti,omitempty"`
	Roles []string        `json:"roles,omitempty"`
	jwt.StandardClaims
}

// JWKResponse ...
type JWKResponse struct {
	Keys []JWKData `json:"keys"`
}

type JWKData struct {
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Key string `json:"k"`
	Kid string `json:"kid"`
}

type JWKGet struct {
	Key       string
	Type      string
	Algorithm string
}
