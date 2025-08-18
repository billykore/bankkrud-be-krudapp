package server

import (
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"go.bankkrud.com/backend/svc/tapmoney/internal/adapter/http/handler"
	"go.bankkrud.com/backend/svc/tapmoney/internal/adapter/log"
	"go.bankkrud.com/backend/svc/tapmoney/internal/pkg/config"
)

// Server represents the main server struct managing configuration, logging, and routing.
type Server struct {
	cfg    *config.Configs
	log    *log.ZapLogger
	router *echo.Echo
	th     *handler.TapMoneyHandler
}

// New returns new Router.
func New(
	cfg *config.Configs,
	log *log.ZapLogger,
	router *echo.Echo,
	th *handler.TapMoneyHandler,
) *Server {
	return &Server{
		cfg:    cfg,
		log:    log,
		router: router,
		th:     th,
	}
}

// Run runs the server router.
func (srv *Server) Run() {
	srv.useMiddlewares()
	srv.registerRoutes()
	srv.run()
}

func (srv *Server) useMiddlewares() {
	srv.router.Use(echomiddleware.Logger())
	srv.router.Use(echomiddleware.Recover())
}

func (srv *Server) run() {
	port := srv.cfg.App.Port
	if err := srv.router.Start(":" + port); err != nil {
		srv.log.Usecase("run").Fatalf("Failed to start server: %v", err)
	}
}
