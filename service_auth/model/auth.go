package model

import (
	"database/sql"
	"strings"
	"time"

	"github.com/syukri21/mercari/common/helper"
)

// RefreshAccessToken ...
type RefreshAccessToken struct {
	AccessToken string
}

// LoginRequest ...
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=0"`
	DeviceID string `json:"deviceId" validate:"required"`
}

// LoginResponse ...
type LoginResponse struct {
	Email             string `json:"email"`
	Username          string `json:"username"`
	ActiveAccessToken string `json:"activeAccessToken"`
}

// RegisterRequest ...
type RegisterRequest struct {
	Email           string `json:"email" from:"email" validate:"required,email"`
	Password        string `json:"password" from:"password" validate:"required,gte=6"`
	ConfirmPassword string `json:"confirmPassword" from:"confirm_password" validate:"required,gte=6"`
}

func (r *RegisterRequest) ToCreateUserRequest() CreateUserRequest {
	return CreateUserRequest{
		Username:    strings.Split(r.Email, "@")[0],
		Email:       r.Email,
		Password:    r.Password,
		ActivateKey: helper.GeneratePin(6),
		IsActive:    false,
		CreatedAt:   time.Now(),
	}
}

// RegisterResponse ...
type RegisterResponse struct {
	Email    string `json:"email" from:"email" `
	Username string `json:"username" `
}

// LogoutRequest ...
type LogoutRequest struct {
}

// LogoutResponse ...
type LogoutResponse struct {
}

type VerifyRegisterRequest struct {
	Email       string `json:"username" validate:"required"`
	ActivateKey string `json:"activateKey" validate:"required"`
}

type VerifyRegisterResponse struct {
}

type LoginHistoryRequest struct {
	Limit  int    `json:"limit" validate:"required" form:"limit"  db:"limit"`
	Offset int    `json:"offset" validate:"required" form:"offset" db:"offset"`
	Email  string `json:"-" db:"email"`
}

type LoginHistoryResponse struct {
	Data []LoginHistory `json:"data"`
}

type LoginHistory struct {
	ID        int          `json:"ID" db:"id"`
	Email     string       `json:"email" db:"email"`
	Username  string       `json:"username" db:"username"`
	DeviceId  string       `json:"deviceId" db:"device_id"`
	LoginAt   time.Time    `json:"loginAt" db:"login_at"`
	CreatedAt time.Time    `json:"CreatedAt" db:"created_at"`
	UpdateAt  sql.NullTime `json:"updateAt" db:"updated_at"`
}
