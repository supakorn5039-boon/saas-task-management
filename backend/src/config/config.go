package config

import (
	"fmt"
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

	jwtSecret := envOr("JWT_SECRET", cfg.Section("server").Key("jwt_secret").String())
	if jwtSecret == "" {
		return fmt.Errorf("jwt_secret is required (set [server] jwt_secret in config.ini or JWT_SECRET env var)")
	}
	if len(jwtSecret) < 32 {
		return fmt.Errorf("jwt_secret must be at least 32 characters (got %d)", len(jwtSecret))
	}

	App = &Config{
		Server: ServerConfig{
			Port:       cfg.Section("server").Key("port").MustInt(8080),
			Production: cfg.Section("server").Key("production").MustBool(false),
			JWTSecret:  jwtSecret,
		},
		Database: DatabaseConfig{
			Host:     envOr("DB_HOST", cfg.Section("database").Key("host").MustString("localhost")),
			Port:     cfg.Section("database").Key("port").MustInt(5432),
			User:     envOr("DB_USER", cfg.Section("database").Key("user").MustString("postgres")),
			Password: envOr("DB_PASSWORD", cfg.Section("database").Key("password").String()),
			Database: envOr("DB_NAME", cfg.Section("database").Key("database").MustString("postgres")),
		},
	}
	return nil
}

func envOr(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
