package postgres

import (
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/config"
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

func Close(db *gorm.DB) error {
	sql, err := db.DB()
	if err != nil {
		return err
	}
	return sql.Close()
}
