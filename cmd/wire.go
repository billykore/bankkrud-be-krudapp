//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"go.bankkrud.com/backend/svc/tapmoney/internal/adapter"
	"go.bankkrud.com/backend/svc/tapmoney/internal/pkg"
	"go.bankkrud.com/backend/svc/tapmoney/internal/pkg/config"
	"go.bankkrud.com/backend/svc/tapmoney/internal/usecase"
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
