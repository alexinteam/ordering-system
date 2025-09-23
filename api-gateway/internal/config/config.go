package config

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server   ServerConfig
	Services ServicesConfig
	Debug    bool `envconfig:"DEBUG" default:"false"`
}

type ServerConfig struct {
	Host string `envconfig:"SERVER_HOST" default:"0.0.0.0"`
	Port int    `envconfig:"SERVER_PORT" default:"8080"`
}

type ServicesConfig struct {
	BillingService      ServiceConfig `envconfig:"BILLING_SERVICE"`
	NotificationService ServiceConfig `envconfig:"NOTIFICATION_SERVICE"`
	OrderService        ServiceConfig `envconfig:"ORDER_SERVICE"`
}

type ServiceConfig struct {
	URL string `envconfig:"URL"`
}

func Load() *Config {
	var config Config

	if err := envconfig.Process("GATEWAY", &config); err != nil {
		log.Fatal("Failed to process environment variables:", err)
	}

	return &config
}

func (c *ServerConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
