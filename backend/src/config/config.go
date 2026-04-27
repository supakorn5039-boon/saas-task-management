package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/ini.v1"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	CORS     CORSConfig
}

type ServerConfig struct {
	Port       int
	Production bool
	JWTSecret  string
}

type DatabaseConfig struct {
	// Either DSN (a full Postgres connection string) is set, or the discrete
	// fields below — DSN wins when present. Cloud providers (Neon, Render,
	// Supabase) hand out a DSN, so DSN is the path of least resistance there.
	DSN string

	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type CORSConfig struct {
	// AllowedOrigins is a comma-separated list of origins. Defaults to the
	// local dev URLs; in production set FRONTEND_URL to e.g.
	// "https://your-app.vercel.app".
	AllowedOrigins []string
}

var App *Config

// Load reads configuration in this order of precedence (highest first):
//  1. Environment variables (JWT_SECRET, DATABASE_URL, PORT, etc.)
//  2. The config.ini at `path` (if present)
//  3. Sensible defaults for local dev
//
// In cloud deploys (Render / Fly / Vercel) there's no config.ini — just
// env vars — and Load is happy with that. The file is only required when
// neither it nor the JWT_SECRET env var is set.
func Load(path string) error {
	cfg := loadINIIfPresent(path)

	jwtSecret := envOr("JWT_SECRET", iniGet(cfg, "server", "jwt_secret", ""))
	if jwtSecret == "" {
		return fmt.Errorf("jwt_secret is required (set JWT_SECRET env var or [server] jwt_secret in config.ini)")
	}
	if len(jwtSecret) < 32 {
		return fmt.Errorf("jwt_secret must be at least 32 characters (got %d)", len(jwtSecret))
	}

	App = &Config{
		Server: ServerConfig{
			// PORT is the cloud-provider convention (Render/Heroku/Fly all
			// inject it). Falls back to the INI value, then 8080.
			Port:       envInt("PORT", iniGetInt(cfg, "server", "port", 8080)),
			Production: envBool("PRODUCTION", iniGetBool(cfg, "server", "production", false)),
			JWTSecret:  jwtSecret,
		},
		Database: DatabaseConfig{
			DSN:      envOr("DATABASE_URL", ""),
			Host:     envOr("DB_HOST", iniGet(cfg, "database", "host", "localhost")),
			Port:     envInt("DB_PORT", iniGetInt(cfg, "database", "port", 5432)),
			User:     envOr("DB_USER", iniGet(cfg, "database", "user", "postgres")),
			Password: envOr("DB_PASSWORD", iniGet(cfg, "database", "password", "")),
			Database: envOr("DB_NAME", iniGet(cfg, "database", "database", "postgres")),
		},
		CORS: CORSConfig{
			AllowedOrigins: parseOrigins(envOr(
				"FRONTEND_URL",
				"http://localhost:5173,http://saas-management.local",
			)),
		},
	}
	return nil
}

// loadINIIfPresent returns nil if config.ini doesn't exist (cloud-deploy case).
// Callers must handle the nil case via iniGet helpers.
func loadINIIfPresent(path string) *ini.File {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		path = "../" + path
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	cfg, err := ini.Load(path)
	if err != nil {
		return nil
	}
	return cfg
}

func iniGet(cfg *ini.File, section, key, fallback string) string {
	if cfg == nil {
		return fallback
	}
	v := cfg.Section(section).Key(key).String()
	if v == "" {
		return fallback
	}
	return v
}

func iniGetInt(cfg *ini.File, section, key string, fallback int) int {
	if cfg == nil {
		return fallback
	}
	return cfg.Section(section).Key(key).MustInt(fallback)
}

func iniGetBool(cfg *ini.File, section, key string, fallback bool) bool {
	if cfg == nil {
		return fallback
	}
	return cfg.Section(section).Key(key).MustBool(fallback)
}

func envOr(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

func envBool(key string, fallback bool) bool {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}

func parseOrigins(raw string) []string {
	out := []string{}
	start := 0
	for i := 0; i <= len(raw); i++ {
		if i == len(raw) || raw[i] == ',' {
			s := trimSpace(raw[start:i])
			if s != "" {
				out = append(out, s)
			}
			start = i + 1
		}
	}
	return out
}

func trimSpace(s string) string {
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\t') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\t') {
		s = s[:len(s)-1]
	}
	return s
}
