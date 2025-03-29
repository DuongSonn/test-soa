package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Configuration struct {
	PostgresDB PostgresDatabase `mapstructure:"postgres"`
	Server     Server           `mapstructure:"server"`
	Jwt        JWT              `mapstructure:"jwt"`
	Redis      Redis            `mapstructure:"redis"`
}

// NewConfigClient creates a new configuration client
func NewConfigClient(filePath string) (*Configuration, error) {
	var configuration Configuration

	v := viper.New()
	v.SetConfigFile(filePath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	if err := v.Unmarshal(&configuration); err != nil {
		return nil, fmt.Errorf("error decoding config: %v", err)
	}

	// Set default values
	if configuration.Server.Port == 0 {
		configuration.Server.Port = 8080
	}

	return &configuration, nil
}

// Get returns the current configuration
func (c *Configuration) Get() Configuration {
	return *c
}
