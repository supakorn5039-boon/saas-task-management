package app

import (
	"github.com/supakorn5039-boon/saas-task-backend/src/config"
	"github.com/supakorn5039-boon/saas-task-backend/src/database"
	"github.com/supakorn5039-boon/saas-task-backend/src/security"
	"github.com/supakorn5039-boon/saas-task-backend/src/server"
)

type App struct {
	config *config.Config
}

func NewApp(cfg *config.Config) *App {
	database.Init(&cfg.Database)
	security.InitJWT(cfg.Server.JWTSecret)
	return &App{config: cfg}
}

func (a *App) WebServer() {
	server.WebServer(&a.config.Server)
}
