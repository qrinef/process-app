package config

import (
	"time"

	commonConfig "github.com/qrinef/process-app/pkg/config"
	commonLogger "github.com/qrinef/process-app/pkg/logger"
)

type Config struct {
	// -------------------
	// Application configs
	// -------------------
	*commonConfig.AppConfig
	*commonLogger.LoggerConfig

	// -------------------
	// Internal configs
	// -------------------
	ClientBufferSize uint64 `envconfig:"CLIENT_BUFFER_SIZE" default:"100"`

	ExternalBatchItemsLimit uint64        `envconfig:"EXTERNAL_BATCH_ITEMS_LIMIT" default:"10"`
	ExternalIntervalLimit   time.Duration `envconfig:"EXTERNAL_INTERVAL_LIMIT" default:"3s"`
}

func (c *Config) GetClientBufferSize() uint64 {
	return c.ClientBufferSize
}

func (c *Config) SetClientBufferSize(value uint64) {
	c.ClientBufferSize = value
}

func (c *Config) GetExternalBatchItemsLimit() uint64 {
	return c.ExternalBatchItemsLimit
}

func (c *Config) GetExternalIntervalLimit() time.Duration {
	return c.ExternalIntervalLimit
}
