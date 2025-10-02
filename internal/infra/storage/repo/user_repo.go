package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/user"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/infra/storage/model"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/config"
)

const (
	userTokenKey        = "user:%s:token"
	duplicateKeyErrCode = "23505"
)

type UserRepo struct {
	cfg *config.Configs
	db  *gorm.DB
	rdb *redis.Client
}

func NewUserRepo(cfg *config.Configs, db *gorm.DB, rdb *redis.Client) *UserRepo {
	return &UserRepo{
		cfg: cfg,
		db:  db,
		rdb: rdb,
	}
}

func (r *UserRepo) Create(ctx context.Context, u user.User) error {
	m := model.User{
		Email:        u.Email,
		Username:     u.Username,
		PasswordHash: u.Password,
		PhoneNumber:  u.PhoneNumber,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		CIF:          u.CIF,
		Address:      u.Address,
		LastLogin:    time.Now(),
		DateOfBirth:  u.DateOfBirth,
		Status:       u.Status,
	}
	err := r.db.WithContext(ctx).Create(&m).Error
	var pgconnErr *pgconn.PgError
	if err != nil && errors.As(err, &pgconnErr) {
		if pgconnErr.Code != duplicateKeyErrCode {
			return err
		}
		k := getDuplicateKey(pgconnErr.Detail)
		if k != "" {
			return fmt.Errorf("%w %s", user.ErrDuplicateUserData, k)
		}
	}
	return nil
}

func getDuplicateKey(detail string) string {
	if detail == "" {
		return ""
	}
	s := strings.Split(detail, " ")
	if len(s) < 2 {
		return ""
	}
	keyValue := strings.Split(s[1], "=")
	if len(keyValue) < 2 {
		return ""
	}
	return keyValue[0]
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (user.User, error) {
	var m model.User
	err := r.db.WithContext(ctx).
		Where("username = ?", username).
		First(&m).Error
	if err != nil {
		return user.User{}, err
	}
	return user.User{
		Email:       m.Email,
		Username:    m.Username,
		Password:    m.PasswordHash,
		PhoneNumber: m.PhoneNumber,
		FirstName:   m.FirstName,
		LastName:    m.LastName,
		CIF:         m.CIF,
		Address:     m.Address,
		LastLogin:   m.LastLogin,
	}, nil
}

func (r *UserRepo) GetFieldsByUsername(ctx context.Context, username string, fields ...string) (user.User, error) {
	var m model.User
	err := r.db.WithContext(ctx).
		Select(fields).
		Where("username = ?", username).
		First(&m).Error
	if err != nil {
		return user.User{}, err
	}
	return user.User{
		Email:       m.Email,
		Username:    m.Username,
		Password:    m.PasswordHash,
		PhoneNumber: m.PhoneNumber,
		FirstName:   m.FirstName,
		LastName:    m.LastName,
		CIF:         m.CIF,
		Address:     m.Address,
		DateOfBirth: m.DateOfBirth,
		LastLogin:   m.LastLogin,
	}, nil
}

func (r *UserRepo) SaveToken(ctx context.Context, username string, token user.Token) error {
	b, err := json.Marshal(model.Token{
		Value:     token.Value,
		ExpiresAt: token.ExpiresAt,
	})
	if err != nil {
		return err
	}
	redisKey := fmt.Sprintf(userTokenKey, username)
	err = r.rdb.Set(ctx, redisKey, string(b), r.cfg.Token.Duration).Err()
	if err != nil {
		return err
	}
	err = r.db.WithContext(ctx).Model(model.User{}).
		Where("username = ?", username).
		UpdateColumn("last_login", time.Now()).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) GetToken(ctx context.Context, username string) (user.Token, error) {
	redisKey := fmt.Sprintf(userTokenKey, username)
	token, err := r.rdb.Get(ctx, redisKey).Result()
	if err != nil && errors.Is(err, redis.Nil) {
		return user.Token{}, user.ErrTokenNotFound
	}
	if err != nil {
		return user.Token{}, err
	}
	var m model.Token
	err = json.Unmarshal([]byte(token), &m)
	if err != nil {
		return user.Token{}, err
	}
	return user.Token{
		Value:     m.Value,
		ExpiresAt: m.ExpiresAt,
	}, nil
}

func (r *UserRepo) DeleteToken(ctx context.Context, username string) error {
	redisKey := fmt.Sprintf(userTokenKey, username)
	return r.rdb.Del(ctx, redisKey).Err()
}
