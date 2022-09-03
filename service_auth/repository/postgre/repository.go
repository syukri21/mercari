package postgre

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/syukri21/mercari/service_auth/constant"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/syukri21/mercari/service_auth/model"
	"github.com/syukri21/mercari/service_auth/repository"
)

// PostgreRepository ...
type PostgreRepository struct {
	l  *log.Logger
	db *sqlx.DB
}

func (p *PostgreRepository) GetLoginHistories(ctx context.Context, req model.LoginHistoryRequest) (result []model.LoginHistory, err error) {
	query := GetLoginHistories

	row, err := p.db.NamedQueryContext(ctx, query, req)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		history := model.LoginHistory{}
		err := row.Scan(
			&history.ID,
			&history.Email,
			&history.Username,
			&history.DeviceId,
			&history.LoginAt,
			&history.CreatedAt,
			&history.UpdateAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, history)
	}
	return result, err
}

func (p *PostgreRepository) CreateLoginHistory(ctx context.Context, req model.LoginHistory) error {
	query := CreateLoginHistory
	_, err := p.db.NamedExecContext(ctx, query, req)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgreRepository) ValidateUser(ctx context.Context, email string, pin string) (err error) {
	queryGetPin := GetUserPinByEmail
	var dbPin string
	err = p.db.GetContext(ctx, &dbPin, queryGetPin, email)
	if err != nil {
		p.l.Printf("[Error when get pin user] %s", err)
		return err
	}

	if dbPin != pin {
		p.l.Printf("[Error pin not match] %s", err)
		return errors.New(constant.StatusUnauthorized)
	}

	_, err = p.db.NamedExecContext(ctx, ActivateUser, map[string]interface{}{
		"updated_at": time.Now(),
		"email":      email,
	})
	if err != nil {
		p.l.Printf("[Error when get pin user] %s", err)
		return err
	}

	return nil
}

func (p *PostgreRepository) CreateUser(ctx context.Context, request model.CreateUserRequest) (err error) {
	query := CreateUser

	_, err = p.db.NamedQueryContext(ctx, query, request)
	if err != nil {
		p.l.Printf("[Error when create user] %s", err)
		return err
	}

	return err
}

func (p *PostgreRepository) GetUserByEmail(ctx context.Context, req model.LoginRequest) (user model.User, err error) {
	query := GetUser

	err = p.db.GetContext(ctx, &user, query, req.Email)
	if err == sql.ErrNoRows {
		return user, fmt.Errorf(constant.StatusUnauthorized)
	}

	if err != nil {
		p.l.Printf("[Error when GetUser] %s", err)
		return user, err
	}

	return
}

// NewPostgreRepository ...
func NewPostgreRepository(l *log.Logger, db *sqlx.DB) repository.PostgreRepository {
	return &PostgreRepository{
		l:  l,
		db: db,
	}
}
