package server

import "go.bankkrud.com/bankkrud/backend/krudapp/internal/adapter/http/middleware"

func (hs *HTTPServer) registerRoutes() {
	v1 := hs.router.Group("/v1")

	v1.POST("/authentication/login", hs.ah.Login)

	withAuth := v1.Group("", middleware.AuthenticateUser(hs.cfg))

	withAuth.POST("/tapmoney/inquiry", hs.tmh.Inquiry)
	withAuth.POST("/tapmoney/payment", hs.tmh.Payment)

	withAuth.POST("/transfer/init", hs.tfh.Initiate)
	withAuth.POST("/transfer/process", hs.tfh.Process)
	withAuth.GET("/transfer/:uuid", hs.tfh.Detail)

	withAuth.GET("/users/me", hs.uh.GetByUsername)
}
