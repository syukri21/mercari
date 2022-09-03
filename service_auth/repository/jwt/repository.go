package jwt

import (
	"encoding/base64"
	"errors"
	"github.com/square/go-jose/v3"
	"github.com/syukri21/mercari/service_auth/constant"
	"log"

	gojwt "github.com/golang-jwt/jwt/v4"
	"github.com/square/go-jose/v3/jwt"

	"github.com/syukri21/mercari/service_auth/model"
	"github.com/syukri21/mercari/service_auth/repository"
)

// JWTRepository ...
type JWTRepository struct {
}

// NewJWTRepository ...
func NewJWTRepository() repository.JWTRepository {
	return &JWTRepository{}
}

func (f *JWTRepository) DecodeJWTToken(token string) (*model.CustomClaimData, error) {
	parser := &gojwt.Parser{}
	rt, _, err := parser.ParseUnverified(token, new(model.PlainToken))
	if err != nil {
		return nil, err
	}

	claim, ok := rt.Claims.(*model.PlainToken)
	if !ok {
		return nil, errors.New(constant.StatusUnauthorized)
	}
	return &claim.Data, nil
}

func (f *JWTRepository) GenerateJWTToken(plainToken model.PlainToken, jwtGet model.JWKGet) (string, error) {
	sDec, err := base64.StdEncoding.DecodeString(jwtGet.Key)
	if err != nil {
		log.Printf("error decoding. err %s", err.Error())
		return "", err
	}

	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jwtGet.Algorithm, Key: sDec}, &jose.SignerOptions{})
	if err != nil {
		return "", err
	}

	raw, err := jwt.Signed(sig).Claims(plainToken).CompactSerialize()
	if err != nil {
		return "", err
	}

	return raw, nil
}
