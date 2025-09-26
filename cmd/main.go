package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/infra/http/server"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/config"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/db/postgres"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/log"
	"gorm.io/gorm"
)

// main swaggo annotation.
//
//	@title			API Specification
//	@version		1.0
//	@description	Bankfrud service API specification.
//	@termsOfService	https://swagger.io/terms/
//	@contact.name	Billy Kore
//	@contact.url	https://www.swagger.io/support
//	@contact.email	billyimmcul2010@gmail.com
//	@license.name	Apache 2.0
//	@license.url	https://www.apache.org/licenses/LICENSE-2.0.html
//	@host			api.bankkrud.com
//	@schemes		http https
//	@BasePath		/v1
func main() {
	c := config.Load()
	a := initKrudApp(c)
	log.Configure(c.App.Env)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// run http server
	go a.http.Run()

	// wait for termination syscalls and doing cleanup operations after received it
	wait := gracefulShutdown(ctx, 3*time.Second, map[string]operation{
		"postgres": func(ctx context.Context) error {
			return postgres.Close(a.db)
		},
		"redis": func(ctx context.Context) error {
			return a.rds.Close()
		},
		"http-server": func(ctx context.Context) error {
			return a.http.Shutdown(ctx)
		},
	})

	<-wait
}

type krudApp struct {
	http *server.HTTPServer
	db   *gorm.DB
	rds  *redis.Client
}

func newKrudApp(http *server.HTTPServer, db *gorm.DB, rds *redis.Client) *krudApp {
	return &krudApp{
		http: http,
		db:   db,
		rds:  rds,
	}
}

// operation is a function that performs some cleanup operation.
type operation func(ctx context.Context) error

// gracefulShutdown waits for termination syscalls and doing cleanup operations after received it
func gracefulShutdown(ctx context.Context, timeout time.Duration, ops map[string]operation) <-chan struct{} {
	l := log.WithContext(ctx, "gracefulShutdown")

	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)
		// add any other syscalls that you want to be notified with
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s
		l.Info().Msg("shutting down")
		// set timeout for the ops to be done to prevent system hang
		timeoutFunc := time.AfterFunc(timeout, func() {
			l.Info().Msgf("timeout %d ms has been elapsed, force exit", timeout.Milliseconds())
			os.Exit(0)
		})
		defer timeoutFunc.Stop()

		var wg sync.WaitGroup
		// do the operations asynchronously to save time
		for key, op := range ops {
			wg.Add(1)
			innerOp := op
			innerKey := key
			go func() {
				defer wg.Done()
				l.Info().Msgf("cleaning up: %s", innerKey)
				if err := innerOp(ctx); err != nil {
					l.Info().Msgf("%s: clean up failed: %v", innerKey, err)
					return
				}
				l.Info().Msgf("%s was shutdown gracefully", innerKey)
			}()
		}
		wg.Wait()
		close(wait)
	}()

	return wait
}
