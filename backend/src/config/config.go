package config

import (
	"os"

	"gopkg.in/ini.v1"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port       int
	Production bool
	JWTSecret  string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

var App *Config

func Load(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		path = "../" + path
	}

	cfg, err := ini.Load(path)
	if err != nil {
		return err
	}

	App = &Config{
		Server: ServerConfig{
			Port:       cfg.Section("server").Key("port").MustInt(8080),
			Production: cfg.Section("server").Key("production").MustBool(false),
			JWTSecret:  cfg.Section("server").Key("jwt_secret").String(),
		},
		Database: DatabaseConfig{
			Host:     cfg.Section("database").Key("host").MustString("localhost"),
			Port:     cfg.Section("database").Key("port").MustInt(5432),
			User:     cfg.Section("database").Key("user").MustString("postgres"),
			Password: cfg.Section("database").Key("password").String(),
			Database: cfg.Section("database").Key("database").MustString("postgres"),
		},
	}
	return nil
}
