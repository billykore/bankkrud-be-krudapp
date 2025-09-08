package server

func (hs *HTTPServer) registerRoutes() {
	v1 := hs.router.Group("/api/v1")

	v1.POST("/tapmoney/inquiry", hs.tmh.Inquiry)
	v1.POST("/tapmoney/payment", hs.tmh.Payment)

	v1.POST("/transfer/initiate", hs.tfh.Initiate)
	v1.POST("/transfer/process", hs.tfh.Process)
}
