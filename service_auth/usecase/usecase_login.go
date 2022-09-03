package usecase

import (
	"context"
	"errors"
	"github.com/square/go-jose/v3"
	"github.com/syukri21/mercari/common/helper"
	"github.com/syukri21/mercari/service_auth/constant"
	"github.com/syukri21/mercari/service_auth/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"time"
)

func (a *AuthUsecase) Login(ctx context.Context, req model.LoginRequest) (result model.LoginResponse, err error) {
	_, span := otel.Tracer(ServicesName).Start(ctx, "Login")
	defer span.End()

	user, err := a.Postgre.GetUserByEmail(ctx, req)
	if err != nil {
		a.Logger.Printf("[Error a.Postgre.GetUserByEmail; %s]", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return model.LoginResponse{}, err
	}

	if !helper.CompareHashAndPassword(user.Password, req.Password) {
		a.Logger.Print("[Not Match CompareHashAndPassword]")
		span.SetStatus(codes.Unset, constant.StatusUnauthorized)
		return model.LoginResponse{}, errors.New(constant.StatusUnauthorized)
	}

	tokenClaim := model.GenerateCustomClaim(model.RequestCustomClaim{
		Username: user.Username,
		DeviceID: req.DeviceID,
		Email:    user.Email,
	})

	at, rt, err := a.generateJWTToken(tokenClaim)
	if err != nil {
		a.Logger.Printf("[generateJWTToken] error; %s", err)
		span.SetStatus(codes.Error, err.Error())
		return model.LoginResponse{}, errors.New(constant.StatusSystemError)
	}

	err = a.Postgre.CreateLoginHistory(ctx, model.LoginHistory{
		Email:     user.Email,
		Username:  user.Username,
		DeviceId:  req.DeviceID,
		LoginAt:   time.Now(),
		CreatedAt: time.Now(),
	})
	if err != nil {
		a.Logger.Printf("[a.Postgre.CreateLoginHistory] error; %s", err)
		span.SetStatus(codes.Error, constant.StatusSystemError)
		return model.LoginResponse{}, errors.New(constant.StatusSystemError)
	}

	err = a.Redis.SaveLoginToken(ctx, rt, at, user.Username, user.Email, req.DeviceID)
	if err != nil {
		a.Logger.Printf("[a.Redis.SaveLoginToken] error; %s", err)
		span.SetStatus(codes.Error, constant.StatusSystemError)
		return model.LoginResponse{}, errors.New(constant.StatusSystemError)
	}

	return model.LoginResponse{
		Email:             user.Email,
		Username:          user.Username,
		ActiveAccessToken: at,
	}, nil
}

func (a *AuthUsecase) RefreshAccessToken(ctx context.Context, deviceId string, email string, activeToken string) (result model.RefreshAccessToken, err error) {
	_, span := otel.Tracer(ServicesName).Start(ctx, "RefreshAccessToken")
	defer span.End()

	userToken, err := a.Redis.GetLoginToken(ctx, email)
	if err != nil {
		return result, errors.New(constant.StatusUnauthorized)
	}

	if userToken.ActiveToken != activeToken {
		return result, errors.New(constant.StatusUnauthorized)
	}

	decoded, err := a.JWT.DecodeJWTToken(userToken.RefreshToken)
	if err != nil {
		return result, errors.New(constant.StatusUnauthorized)
	}

	tokenClaim := model.GenerateCustomClaim(model.RequestCustomClaim{
		Username: decoded.Username,
		DeviceID: decoded.DeviceID,
		Email:    decoded.Email,
	})

	at, rt, err := a.generateJWTToken(tokenClaim)
	if err != nil {
		a.Logger.Printf("[generateJWTToken] error; %s", err)
		span.SetStatus(codes.Error, constant.StatusSystemError)
		return model.RefreshAccessToken{}, errors.New(constant.StatusSystemError)
	}

	err = a.Redis.SaveLoginToken(ctx, rt, at, decoded.Username, decoded.Email, deviceId)
	if err != nil {
		a.Logger.Printf("[generateJWTToken] error; %s", err)
		span.SetStatus(codes.Error, constant.StatusSystemError)
		return result, errors.New(constant.StatusSystemError)
	}

	result.AccessToken = at
	return result, nil
}

func (a *AuthUsecase) Logout(ctx context.Context, email string, activeToken string) error {
	_, span := otel.Tracer(ServicesName).Start(ctx, "Logout")
	defer span.End()

	userToken, err := a.Redis.GetLoginToken(ctx, email)
	if err != nil {
		return errors.New(constant.StatusUnauthorized)
	}

	if userToken.ActiveToken != activeToken {
		return errors.New(constant.StatusUnauthorized)
	}

	err = a.Redis.ClearLoginToken(ctx, email)
	if err != nil {
		a.Logger.Printf("[ClearLoginToken] error; %s", err)
		span.SetStatus(codes.Error, constant.StatusSystemError)
		return err
	}
	return nil
}

func (a *AuthUsecase) GetLoginHistories(ctx context.Context, req model.LoginHistoryRequest) (result model.LoginHistoryResponse, err error) {
	_, span := otel.Tracer(ServicesName).Start(ctx, "GetLoginHistories")
	defer span.End()
	session := ctx.Value(CtxSession).(model.CustomClaimData)
	span.SetAttributes(attribute.String("email", session.Email))
	span.SetAttributes(attribute.String("device_id", session.DeviceID))
	span.SetAttributes(attribute.String("username", session.Username))

	req.Email = session.Email

	histories, err := a.Postgre.GetLoginHistories(ctx, req)
	if err != nil {
		a.Logger.Printf("[a.Redis.GetLoginHistories] error; %s", err)
		span.SetStatus(codes.Error, err.Error())
		return result, errors.New(constant.StatusSystemError)
	}

	return model.LoginHistoryResponse{
		Data: histories,
	}, nil
}

func (a *AuthUsecase) generateJWTToken(claim *model.CustomClaim) (at string, rt string, err error) {
	accessToken := ""
	refreshToken := ""

	plainToken := model.PlainToken{
		Aud:   "https://mercari.com",
		Data:  claim.Data,
		Exp:   int32(claim.ExpiresAt),
		Iss:   "MERCARI",
		Jti:   claim.Id,
		Roles: nil,
	}
	refreshPlainToken := model.PlainToken{
		Aud:   "https://mercari.com",
		Data:  claim.Data,
		Exp:   int32(claim.RefreshExpiresAt),
		Iss:   "MERCARI_TOKEN",
		Jti:   claim.Id,
		Roles: nil,
	}

	jwkGet := model.JWKGet{
		Key:       a.config.App.JWTKey,
		Algorithm: jose.HS256,
	}

	accessToken, err = a.JWT.GenerateJWTToken(plainToken, jwkGet)
	if err != nil {
		return accessToken, refreshToken, errors.New("[SBE 049] Internal error occurred. Please contact our administrator")
	}

	refreshToken, err = a.JWT.GenerateJWTToken(refreshPlainToken, jwkGet)
	if err != nil {
		return accessToken, refreshToken, errors.New("[SBE 140] Internal error occurred. Please contact our administrator")
	}

	return accessToken, refreshToken, nil
}
