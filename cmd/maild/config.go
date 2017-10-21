package main

import (
	"gopkg.in/go-playground/validator.v9"
)

type Config struct {
	BrokerURL    string `validated:"required" default:amqp://localhost:5672`
	ServiceName  string `validated:"required" default:"maild"`
	ProviderURL  string `validated:"required"`
	PrividerKey  string `validated:"required"`
	DatabaseType string `validated:"required"`
	DatabaseDSN  string `validated:"required"`
	LogLevel     string `validated:"required" default:"INFO"`
}

func (c *Config) Validate() error {
	return validator.New().Struct(c)
}
