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
)

const (
	userTokenKey      = "user:%s:token"
	userTokenCacheTTL = time.Hour
)

type UserRepo struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewUserRepo(db *gorm.DB, rdb *redis.Client) *UserRepo {
	return &UserRepo{db: db, rdb: rdb}
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
	}, nil
}

func (ur *UserRepo) SaveToken(ctx context.Context, username string, token user.Token) error {
	redisKey := fmt.Sprintf(userTokenKey, username)
	b, err := json.Marshal(token)
	if err != nil {
		return err
	}
	err = ur.rdb.Set(ctx, redisKey, string(b), userTokenCacheTTL).Err()
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
	var t user.Token
	err = json.Unmarshal([]byte(token), &t)
	if err != nil {
		return user.Token{}, err
	}
	return t, nil
}
