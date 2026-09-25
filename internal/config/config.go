package config

import (
	"fmt"
	"net"
	"os"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Port string
	Host string
}

type LogConfig struct {
	Level string
}

type DBConfig struct {
	Path string
}

type Config struct {
	Server ServerConfig
	DB     DBConfig
	Log    LogConfig
}

func (c *ServerConfig) BaseURL() string {
	return fmt.Sprintf("http://%s:%s/", c.Host, c.Port)
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("no .env file")
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Host: getEnv("HOST", "localhost"),
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "INFO"),
		},
		DB: DBConfig{
			Path: getEnv("DB_PATH", "url-shortener.db"),
		},
	}, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return fallback
}

func (s *ServerConfig) Addr() string {
	return net.JoinHostPort(s.Host, s.Port)
}
