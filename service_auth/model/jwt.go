package model

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"github.com/square/go-jose/v3"
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
	Algorithm jose.SignatureAlgorithm
}

type RequestCustomClaim struct {
	Username string
	DeviceID string
	Email    string
}

var Now = time.Now
var GenerateNewToken = NewTokenID

// GenerateCustomClaim returns custom claim of specific user.
// deviceID is the device id used to access the API and issuer is the token issuer.
func GenerateCustomClaim(requestCustomClaim RequestCustomClaim) *CustomClaim {
	issuedTime := Now()
	expiredTime := issuedTime.Add(time.Minute * time.Duration(viper.GetInt("SESSION_EXPIRED_MINUTE")))
	refreshExpiredTime := issuedTime.AddDate(0, 0, 7)

	claim := new(CustomClaim)
	claim.Id = GenerateNewToken().ValueBase64
	claim.IssuedTime = &issuedTime
	claim.IssuedAt = issuedTime.Unix()
	claim.NotBefore = issuedTime.Unix()
	claim.ExpiresAt = expiredTime.Unix()
	claim.RefreshExpiresAt = refreshExpiredTime.Unix()
	claim.Data.Username = requestCustomClaim.Username
	claim.Data.DeviceID = requestCustomClaim.DeviceID
	claim.Data.Email = requestCustomClaim.Email

	return claim
}

type TokenID struct {
	Value       []byte
	ValueBase64 string
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func NewTokenID() *TokenID {
	v, _ := generateRandomBytes(8)
	v64 := base64.StdEncoding.EncodeToString(v)
	return &TokenID{v, v64}
}
