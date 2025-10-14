package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/infra/http/handler"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/infra/http/middleware"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/config"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/trace"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

// HTTPServer represents the main server struct managing configuration, logging, and routing.
type HTTPServer struct {
	cfg    *config.Configs
	router *echo.Echo
	tracer *trace.Tracer
	tmh    *handler.TapMoneyHandler
	tfh    *handler.TransferHandler
	ah     *handler.AuthenticationHandler
	uh     *handler.UserHandler
	txh    *handler.TransactionHandler
}

// NewHTTP returns new Router.
func NewHTTP(
	cfg *config.Configs,
	router *echo.Echo,
	tracer *trace.Tracer,
	tmh *handler.TapMoneyHandler,
	tfh *handler.TransferHandler,
	ah *handler.AuthenticationHandler,
	uh *handler.UserHandler,
	txh *handler.TransactionHandler,
) *HTTPServer {
	return &HTTPServer{
		cfg:    cfg,
		router: router,
		tracer: tracer,
		tmh:    tmh,
		tfh:    tfh,
		ah:     ah,
		uh:     uh,
		txh:    txh,
	}
}

// Run runs the server router.
func (hs *HTTPServer) Run() {
	hs.setupRouter()
	hs.useMiddlewares()
	hs.registerRoutes()
	hs.run()
}

func (hs *HTTPServer) setupRouter() {
	hs.router.HideBanner = true
}

func (hs *HTTPServer) useMiddlewares() {
	hs.router.Use(echomiddleware.Logger())
	hs.router.Use(echomiddleware.Recover())
	hs.router.Use(otelecho.Middleware(hs.cfg.App.Name))
	hs.router.Use(middleware.DumpBodyWithTracer(hs.tracer.Tracer()))
}

func (hs *HTTPServer) run() {
	port := hs.cfg.App.Port
	err := hs.router.Start(":" + port)
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		log.Info().Msgf("http server closed")
	} else {
		log.Panic().Msgf("Failed to start server: %v", err)
	}
}

func (hs *HTTPServer) Shutdown(ctx context.Context) error {
	return hs.router.Shutdown(ctx)
}
