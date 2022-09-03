package usecase

import (
	"context"
	"errors"
	gojwt "github.com/golang-jwt/jwt/v4"
	"github.com/syukri21/mercari/service_auth/constant"
	"github.com/syukri21/mercari/service_auth/model"
	"google.golang.org/grpc/metadata"
	"strings"
)

func (a *AuthUsecase) CheckSession(ctx context.Context) (result model.CustomClaimData, err error) {
	parser := &gojwt.Parser{}
	jwt := a.getMetadata(ctx)
	token, _, err := parser.ParseUnverified(jwt, new(model.PlainToken))
	if err != nil {
		return result, err
	}

	claim, ok := token.Claims.(*model.PlainToken)
	if !ok {
		return result, errors.New(constant.StatusUnauthorized)
	}

	loginToken, err := a.Redis.GetLoginToken(ctx, claim.Data.Email)
	if err != nil {
		return result, errors.New(constant.StatusUnauthorized)
	}

	if loginToken.ActiveToken != jwt {
		return result, errors.New(constant.StatusUnauthorized)
	}

	return claim.Data, nil
}

// getMetadata ...
func (a *AuthUsecase) getMetadata(ctx context.Context) string {
	jwt := ""

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		authToken := md.Get("authorization")
		if len(authToken) > 0 {
			jwt = strings.Replace(authToken[0], "Bearer ", "", 1)
		}
	}

	return jwt
}
