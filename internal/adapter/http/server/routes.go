package server

func (hs *HTTPServer) registerRoutes() {
	v1 := hs.router.Group("/api/v1")
	v1.POST("/tapmoney/inquiry", hs.th.Inquiry)
}
