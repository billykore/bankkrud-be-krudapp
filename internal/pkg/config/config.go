// Package config contains all the service configuration values.
// The configuration is from the config.yaml file.
package config

import (
	"github.com/spf13/viper"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/config/internal"
)

// Configs hold the application configurations.
type Configs struct {
	// App defines the application configuration.
	App internal.App
	// Postgres defines the postgres database configuration.
	Postgres internal.Postgres
	// Token defines the token configuration.
	Token internal.Token
	// CBS defines the core banking system configuration.
	CBS internal.CBS
	// DBD defines the digital banking delivery system configuration.
	DBD internal.DBD
	// Redis defines the redis database configuration.
	Redis internal.Redis
}

// Config holds the application configuration.
type Config struct {
	Name    string
	Version string
	Configs Configs
}

// Load loads application configuration from a YAML file using Viper.
func Load() *Configs {
	viper.SetConfigName("configs")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	cfg := new(Config)
	err := viper.Unmarshal(cfg)
	if err != nil {
		panic(err)
	}

	return &cfg.Configs
}
