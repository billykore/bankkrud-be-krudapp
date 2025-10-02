package server

import "go.bankkrud.com/bankkrud/backend/krudapp/internal/infra/http/middleware"

func (hs *HTTPServer) registerRoutes() {
	v1 := hs.router.Group("/v1")

	v1.POST("/auth/login", hs.ah.Login)

	v1.POST("/users", hs.uh.Create)

	withAuth := v1.Group("", middleware.AuthorizeUser(hs.cfg))

	withAuth.POST("/tapmoney/init", hs.tmh.Initiate)
	withAuth.POST("/tapmoney/:uuid/process", hs.tmh.Process)

	withAuth.POST("/transfers/init", hs.tfh.Initiate)
	withAuth.POST("/transfers/:uuid/process", hs.tfh.Process)

	withAuth.GET("/transactions", hs.txh.GetTransactions)
	withAuth.GET("/transactions/:uuid", hs.txh.GetTransaction)

	withAuth.GET("/users/me", hs.uh.GetByUsername)
}
