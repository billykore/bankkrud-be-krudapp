package main

import (
	"go.bankkrud.com/backend/svc/tapmoney/internal/adapter/http/server"
	"go.bankkrud.com/backend/svc/tapmoney/internal/pkg/config"
)

type tapmoney struct {
	srv *server.Server
}

func newTapMoney(srv *server.Server) *tapmoney {
	return &tapmoney{
		srv: srv,
	}
}

// main swaggo annotation.
//
//	@title			API Specification
//	@version		1.0
//	@description	TapMoney service API specification.
//	@termsOfService	https://swagger.io/terms/
//	@contact.name	Billy Kore
//	@contact.url	https://www.swagger.io/support
//	@contact.email	billyimmcul2010@gmail.com
//	@license.name	Apache 2.0
//	@license.url	https://www.apache.org/licenses/LICENSE-2.0.html
//	@host			api.bankkrud.com
//	@schemes		http https
//	@BasePath		/api/v1
func main() {
	c := config.Load()
	a := initTapMoney(c)
	a.srv.Run()
}
