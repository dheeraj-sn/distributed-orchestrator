package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	v := viper.New()

	// Defaults for scheduler
	v.SetDefault("scheduler.host", "0.0.0.0:50051")
	v.SetDefault("scheduler.metrics_port", 9090)

	// Defaults for worker
	v.SetDefault("worker.host", "0.0.0.0:50052")
	v.SetDefault("worker.metrics_port", 9091)
	v.SetDefault("worker.concurrency", 4)
	v.SetDefault("worker.worker_id", "worker-default")

	// Defaults for client
	v.SetDefault("client.scheduler_addr", "localhost:50051")

	// Retry defaults
	v.SetDefault("retry.max_attempts", 3)
	v.SetDefault("retry.initial_backoff", "1s")
	v.SetDefault("retry.max_backoff", "1m")

	// Logging defaults
	v.SetDefault("logging.level", "debug")
	v.SetDefault("logging.format", "json")

	// Storage defaults
	v.SetDefault("storage.backend", "memory")

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv() // enable env var overrides, e.g. WORKER_CONCURRENCY

	// Read config file if set
	configPath := v.GetString("CONFIG_PATH")
	if configPath != "" {
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Parse duration strings
	var err error
	cfg.Retry.InitialBackoff, err = time.ParseDuration(v.GetString("retry.initial_backoff"))
	if err != nil {
		return nil, fmt.Errorf("invalid retry.initial_backoff: %w", err)
	}
	cfg.Retry.MaxBackoff, err = time.ParseDuration(v.GetString("retry.max_backoff"))
	if err != nil {
		return nil, fmt.Errorf("invalid retry.max_backoff: %w", err)
	}

	return &cfg, nil
}
