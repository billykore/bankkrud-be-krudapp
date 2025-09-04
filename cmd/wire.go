//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/adapter"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/config"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/usecase"
)

func initTapMoney(cfg *config.Configs) *tapmoney {
	wire.Build(
		adapter.ProviderSet,
		pkg.ProviderSet,
		usecase.ProviderSet,
		echo.New,
		newTapMoney,
	)
	return &tapmoney{}
}
