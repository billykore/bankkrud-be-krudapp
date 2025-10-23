package server

import (
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/infra/http/middleware"
)

// Register all routes.
func (hs *HTTPServer) registerRoutes() {
	hs.authRoutes()
	hs.userRoutes()
	hs.tapMoneyRoutes()
	hs.transferRoutes()
	hs.transactionRoutes()
	hs.scheduleRoutes()
}

// Define authentication routes.
func (hs *HTTPServer) authRoutes() {
	ar := hs.router.Group("/v1/auth")
	ar.POST("/login", hs.ah.Login)
}

// Define user routes.
func (hs *HTTPServer) userRoutes() {
	r := hs.router.Group("/v1/users")
	r.POST("", hs.uh.Create)
	r.GET("/me", hs.uh.GetByUsername, middleware.AuthorizeUser(hs.cfg))
}

// Define TapMoney routes.
func (hs *HTTPServer) tapMoneyRoutes() {
	r := hs.router.Group("/v1/tapmoney")
	r.Use(middleware.AuthorizeUser(hs.cfg))

	r.POST("/init", hs.tmh.Initiate)
	r.POST("/:uuid/process", hs.tmh.Process)
}

// Define transfer routes.
func (hs *HTTPServer) transferRoutes() {
	r := hs.router.Group("/v1/transfers")
	r.Use(middleware.AuthorizeUser(hs.cfg))

	r.POST("/init", hs.tfh.Initiate)
	r.POST("/:uuid/process", hs.tfh.Process)
}

// Define transaction routes.
func (hs *HTTPServer) transactionRoutes() {
	r := hs.router.Group("/v1/transactions")
	r.Use(middleware.AuthorizeUser(hs.cfg))

	r.GET("", hs.txh.GetTransactions)
	r.GET("/:uuid", hs.txh.GetTransaction)
}

// Define schedule routes.
func (hs *HTTPServer) scheduleRoutes() {
	r := hs.router.Group("/v1/schedules")
	r.Use(middleware.AuthorizeUser(hs.cfg))

	r.GET("", hs.sh.GetSchedules)
	r.GET("/:uuid", hs.sh.GetSchedule)
	r.POST("", hs.sh.CreateSchedule)
	r.PATCH("/:uuid", hs.sh.UpdateScheduleStatus)
	r.DELETE("/:uuid", hs.sh.DeleteSchedule)
}
