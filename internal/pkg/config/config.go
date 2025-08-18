// Package config contains all the service configuration values.
// The configuration is from the config.yaml file.
package config

import (
	"github.com/spf13/viper"
	"go.bankkrud.com/backend/svc/tapmoney/internal/pkg/config/internal"
)

// Configs hold the application configurations.
type Configs struct {
	App      internal.App
	Postgres internal.Postgres
	Token    internal.Token
	DBD      internal.DBD
}

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
