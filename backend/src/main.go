package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/supakorn5039-boon/saas-task-backend/src/config"
	"github.com/supakorn5039-boon/saas-task-backend/src/database"
	"github.com/supakorn5039-boon/saas-task-backend/src/database/migration"
	"github.com/supakorn5039-boon/saas-task-backend/src/database/seeder"
	"github.com/supakorn5039-boon/saas-task-backend/src/pkg"
)

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

	r := gin.Default()
	pkg.MountAPIWebServer(r)

	log.Printf("starting http server on :%d", config.App.Server.Port)
	if err := r.Run(fmt.Sprintf(":%d", config.App.Server.Port)); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
