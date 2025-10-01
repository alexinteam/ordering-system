package config

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Debug    bool `envconfig:"DEBUG" default:"false"`
}

type DatabaseConfig struct {
	Host     string `envconfig:"DB_HOST" default:"postgres-service"`
	Port     int    `envconfig:"DB_PORT" default:"5432"`
	User     string `envconfig:"DB_USER" default:"postgres"`
	Password string `envconfig:"DB_PASSWORD" default:"postgres"`
	Name     string `envconfig:"DB_NAME" default:"delivery"`
	SSLMode  string `envconfig:"DB_SSL_MODE" default:"disable"`
}

type ServerConfig struct {
	Host string `envconfig:"SERVER_HOST" default:"0.0.0.0"`
	Port int    `envconfig:"SERVER_PORT" default:"8084"`
}

func Load() *Config {
	var config Config

	if err := envconfig.Process("", &config); err != nil {
		log.Fatal("Failed to process environment variables:", err)
	}

	return &config
}

func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}

func (c *ServerConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
