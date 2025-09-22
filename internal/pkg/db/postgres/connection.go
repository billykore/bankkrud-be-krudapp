package postgres

import (
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// New returns new postgres db connection.
func New(cfg *config.Configs) *gorm.DB {
	dsn := cfg.Postgres.DSN
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if cfg.App.Env == "production" {
		db.Logger = db.Logger.LogMode(logger.Silent)
	}
	return db
}

func Close(db *gorm.DB) error {
	sql, err := db.DB()
	if err != nil {
		return err
	}
	return sql.Close()
}
