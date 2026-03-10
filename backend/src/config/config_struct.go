package config

type ServerConfig struct {
	Port       int    `ini:"port"`
	Production bool   `ini:"production"`
	JWTSecret  string `ini:"jwt_secret"`
}

type DatabaseConfig struct {
	Host     string `ini:"host"`
	User     string `ini:"user"`
	Password string `ini:"password"`
	Database string `ini:"database"`
	Port     int    `ini:"port"`
}

type Config struct {
	Server   ServerConfig   `ini:"server"`
	Database DatabaseConfig `ini:"database"`
}
