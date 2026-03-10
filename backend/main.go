package main

import (
	"log"

	"github.com/supakorn5039-boon/saas-task-backend/src/app"
	"github.com/supakorn5039-boon/saas-task-backend/src/config"
)

func main() {
	appConfig := config.NewAppConfig()
	if err := appConfig.Load("config.ini"); err != nil {
		log.Fatal(err)
	}

	a := app.NewApp(appConfig.Config)
	a.WebServer()
}
