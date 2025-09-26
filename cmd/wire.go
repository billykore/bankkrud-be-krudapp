//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/infra"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/config"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/usecase"
)

func initKrudApp(cfg *config.Configs) *krudApp {
	wire.Build(
		infra.ProviderSet,
		pkg.ProviderSet,
		usecase.ProviderSet,
		echo.New,
		newKrudApp,
	)
	return &krudApp{}
}
