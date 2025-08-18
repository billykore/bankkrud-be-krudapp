package postgres

import (
	"go.bankkrud.com/backend/svc/tapmoney/internal/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// New returns new postgres db connection.
func New(cfg *config.Configs) *gorm.DB {
	dsn := cfg.Postgres.DSN
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
