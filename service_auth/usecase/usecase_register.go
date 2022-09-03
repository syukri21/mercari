package usecase

import (
	"context"
	"errors"
	"github.com/syukri21/mercari/common/helper"
	"github.com/syukri21/mercari/service_auth/constant"
	"github.com/syukri21/mercari/service_auth/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (a *AuthUsecase) Register(ctx context.Context, req model.RegisterRequest) (result model.RegisterResponse, err error) {
	_, span := otel.Tracer(ServicesName).Start(ctx, "Register")
	defer span.End()

	user := req.ToCreateUserRequest()
	user.Password, err = helper.HashPassword(user.Password)
	if err != nil {
		a.Logger.Printf("[Error helper.HashPassword; %s]", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return result, errors.New(constant.StatusSystemError)
	}

	err = a.Postgre.CreateUser(ctx, user)
	if err != nil {
		a.Logger.Printf("[Error when create user; %s]", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return result, errors.New(constant.StatusSystemError)
	}

	span.SetAttributes(attribute.String("request.Email", user.Email))
	span.SetAttributes(attribute.String("request.Username", user.Username))
	return model.RegisterResponse{Email: user.Email, Username: user.Username}, nil
}

func (a *AuthUsecase) VerifyRegister(ctx context.Context, req model.VerifyRegisterRequest) (result model.VerifyRegisterResponse, err error) {
	_, span := otel.Tracer(ServicesName).Start(ctx, "VerifyRegister")
	defer span.End()

	err = a.Postgre.ValidateUser(ctx, req.Email, req.ActivateKey)
	if err != nil {
		a.Logger.Printf("[Error a.Postgre.ValidateUser; %s]", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return result, errors.New(constant.StatusSystemError)
	}

	span.SetAttributes(attribute.String("request.Email", req.Email))
	return result, err
}
