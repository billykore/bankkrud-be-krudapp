package server

import (
	"context"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/adapter/http/handler"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/config"
)

// HTTPServer represents the main server struct managing configuration, logging, and routing.
type HTTPServer struct {
	cfg    *config.Configs
	router *echo.Echo
	tmh    *handler.TapMoneyHandler
	tfh    *handler.TransferHandler
}

// NewHTTP returns new Router.
func NewHTTP(
	cfg *config.Configs,
	router *echo.Echo,
	tmh *handler.TapMoneyHandler,
	tfh *handler.TransferHandler,
) *HTTPServer {
	return &HTTPServer{
		cfg:    cfg,
		router: router,
		tmh:    tmh,
		tfh:    tfh,
	}
}

// Run runs the server router.
func (hs *HTTPServer) Run() {
	hs.useMiddlewares()
	hs.registerRoutes()
	hs.run()
}

func (hs *HTTPServer) useMiddlewares() {
	hs.router.Use(echomiddleware.Logger())
	hs.router.Use(echomiddleware.Recover())
}

func (hs *HTTPServer) run() {
	port := hs.cfg.App.Port
	if err := hs.router.Start(":" + port); err != nil {
		log.Panic().Err(err).Msg("Failed to start server")
	}
}

func (hs *HTTPServer) Shutdown(ctx context.Context) error {
	return hs.router.Shutdown(ctx)
}
