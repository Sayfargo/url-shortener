package core_config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	GRPC struct {
		Port    int           `envconfig:"GRPC_PORT" default:"44044"`
		Timeout time.Duration `envconfig:"GRPC_TIMEOUT" default:"10s"`
	}

	Postgres struct {
		URL               string        `envconfig:"POSTGRES_URL" required:"true"`
		MaxConns          int32         `envconfig:"POSTGRES_MAX_CONNS" default:"16"`
		MinConns          int32         `envconfig:"POSTGRES_MIN_CONNS" default:"4"`
		HealthCheckPeriod time.Duration `envconfig:"POSTGRES_HEALTHCHECK_PERIOD" default:"1m"`
	}

	Redis struct {
		Addr string `envconfig:"REDIS_ADDR" required:"true"`
	}

	Logger struct {
		Level string `envconfig:"LOGGER_LEVEL" default:"debug"`
		Dir   string `envconfig:"LOGGER_DIR" required:"true"`
	}
}

// Load config with error
func Load() (*Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to process envconfig : %w", err)
	}

	return &cfg, nil
}

// Load config with panic
func Mustload() *Config {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		panic(err)
	}

	return &cfg
}
