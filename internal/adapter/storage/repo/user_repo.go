package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"go.bankkrud.com/bankkrud/backend/krudapp/internal/adapter/storage/model"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/domain/user"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/config"
)

const (
	userTokenKey      = "user:%s:token"
	userTokenCacheTTL = time.Hour
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

func (ur *UserRepo) GetByUsername(ctx context.Context, username string) (user.User, error) {
	var m model.User
	err := ur.db.WithContext(ctx).
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

func (ur *UserRepo) SaveToken(ctx context.Context, username string, token user.Token) error {
	b, err := json.Marshal(model.Token{
		Value:     token.Value,
		ExpiresAt: token.ExpiresAt,
	})
	if err != nil {
		return err
	}
	redisKey := fmt.Sprintf(userTokenKey, username)
	err = ur.rdb.Set(ctx, redisKey, string(b), ur.cfg.Token.Duration).Err()
	if err != nil {
		return err
	}
	err = ur.db.WithContext(ctx).Model(model.User{}).
		Where("username = ?", username).
		UpdateColumn("last_login", time.Now()).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepo) GetToken(ctx context.Context, username string) (user.Token, error) {
	redisKey := fmt.Sprintf(userTokenKey, username)
	token, err := ur.rdb.Get(ctx, redisKey).Result()
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
