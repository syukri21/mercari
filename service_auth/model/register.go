package model

import (
	"database/sql"
	"time"
)

type CreateUserRequest struct {
	Username    string       `json:"username" db:"username"`
	Email       string       `json:"email" db:"email"`
	Password    string       `json:"password" db:"password"`
	ActivateKey string       `json:"activateKey" db:"activate_key"`
	IsActive    bool         `json:"isActive" db:"is_activated"`
	CreatedAt   time.Time    `json:"createdAt" db:"created_at"`
	UpdateAt    sql.NullTime `json:"updateAt" db:"updated_at"`
}
