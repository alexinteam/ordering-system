package config

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Services ServicesConfig
	Debug    bool `envconfig:"DEBUG" default:"false"`
}

type DatabaseConfig struct {
	Host     string `envconfig:"DB_HOST" default:"localhost"`
	Port     int    `envconfig:"DB_PORT" default:"5432"`
	User     string `envconfig:"DB_USER" default:"postgres"`
	Password string `envconfig:"DB_PASSWORD" default:"postgres"`
	Name     string `envconfig:"DB_NAME" default:"orders"`
	SSLMode  string `envconfig:"DB_SSL_MODE" default:"disable"`
}

type ServerConfig struct {
	Host string `envconfig:"SERVER_HOST" default:"0.0.0.0"`
	Port int    `envconfig:"SERVER_PORT" default:"8081"`
}

type ServicesConfig struct {
	BillingServiceURL      string `envconfig:"BILLING_SERVICE_URL" default:"http://billing-service:8082"`
	WarehouseServiceURL    string `envconfig:"WAREHOUSE_SERVICE_URL" default:"http://warehouse-service:8083"`
	DeliveryServiceURL     string `envconfig:"DELIVERY_SERVICE_URL" default:"http://delivery-service:8084"`
	NotificationServiceURL string `envconfig:"NOTIFICATION_SERVICE_URL" default:"http://notification-service:8085"`
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
