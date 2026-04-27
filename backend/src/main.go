package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/supakorn5039-boon/saas-task-backend/src/config"
	"github.com/supakorn5039-boon/saas-task-backend/src/database"
	"github.com/supakorn5039-boon/saas-task-backend/src/database/migration"
	"github.com/supakorn5039-boon/saas-task-backend/src/database/seeder"
	"github.com/supakorn5039-boon/saas-task-backend/src/pkg"
)

const shutdownTimeout = 15 * time.Second

func main() {
	if err := config.Load("config.ini"); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := database.Connect(); err != nil {
		log.Fatalf("database error: %v", err)
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "seed":
			if err := seeder.Seed(database.DB); err != nil {
				log.Fatalf("seed error: %v", err)
			}
		case "migrate:status":
			migration.Status(database.DB)
		case "migrate:rollback":
			if err := migration.Rollback(database.DB); err != nil {
				log.Fatalf("rollback error: %v", err)
			}
		default:
			fmt.Println("Available commands:")
			fmt.Println("  (none)              Start the HTTP server")
			fmt.Println("  seed                Run database seeders")
			fmt.Println("  migrate:status      Show migration status")
			fmt.Println("  migrate:rollback    Rollback last migration")
		}
		return
	}

	if config.App.Server.Production {
		gin.SetMode(gin.ReleaseMode)
	}

	// gin.New() (not gin.Default()) so we can plug in our own structured request
	// logger via pkg.MountAPIWebServer; we still want the panic-recovery middleware.
	r := gin.New()
	r.Use(gin.Recovery())
	pkg.MountAPIWebServer(r)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", config.App.Server.Port),
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Run the server in a goroutine so the main goroutine can wait for SIGTERM
	// and trigger a graceful shutdown — drains in-flight requests instead of
	// dropping them when the pod gets killed.
	go func() {
		log.Printf("starting http server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop
	log.Printf("received %s, shutting down...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
	log.Println("server stopped cleanly")
}
