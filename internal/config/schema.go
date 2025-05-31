package config

import (
	"time"
)

type SchedulerConfig struct {
	Host        string `mapstructure:"host"`
	MetricsPort int    `mapstructure:"metrics_port"`
}

type WorkerConfig struct {
	Host        string `mapstructure:"host"`
	MetricsPort int    `mapstructure:"metrics_port"`
	Concurrency int    `mapstructure:"concurrency"`
	WorkerID    string `mapstructure:"worker_id"`
}

type ClientConfig struct {
	SchedulerAddr string `mapstructure:"scheduler_addr"`
}

type RetryConfig struct {
	MaxAttempts    int           `mapstructure:"max_attempts"`
	InitialBackoff time.Duration `mapstructure:"initial_backoff"`
	MaxBackoff     time.Duration `mapstructure:"max_backoff"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type StorageConfig struct {
	Backend string `mapstructure:"backend"`
}

type Config struct {
	Scheduler SchedulerConfig `mapstructure:"scheduler"`
	Worker    WorkerConfig    `mapstructure:"worker"`
	Client    ClientConfig    `mapstructure:"client"`
	Retry     RetryConfig     `mapstructure:"retry"`
	Logging   LoggingConfig   `mapstructure:"logging"`
	Storage   StorageConfig   `mapstructure:"storage"`
}
