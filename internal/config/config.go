package config

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"thumburl-service/internal/pkg/logger"
)

type PoolConfig struct {
	URL          string `json:"url"`
	Count        int    `json:"count"`
	MaxUsedTimes int    `json:"max_used_times"`
	TimeoutSec   int    `json:"timeout_sec"`
}

type ConfigType struct {
	Port               string        `json:"port"`
	PoolConfig         []*PoolConfig `json:"pool_config"`
	GetAgentTimeoutSec int           `json:"get_agent_timeout_sec"`
}

var Config = &ConfigType{
	Port: ":8080",
	PoolConfig: []*PoolConfig{
		{
			URL:          "http://localhost:9222",
			Count:        10,
			MaxUsedTimes: 20,
			TimeoutSec:   10,
		},
	},
	GetAgentTimeoutSec: 5,
}

func Init() {
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		file, err := os.Open(path)
		if err != nil {
			logger.Panicw(
				context.Background(),
				"open file error",
				"err", err,
			)
		}
		defer file.Close()
		content, err := io.ReadAll(file)
		if err != nil {
			logger.Panicw(
				context.Background(),
				"reading file error",
				"err", err,
			)
		}
		if err := json.Unmarshal(content, Config); err != nil {
			logger.Panicw(
				context.Background(),
				"json unmarshal error",
				"err", err,
			)
		}
		logger.Infow(
			context.Background(),
			"reading config from file",
			"path", path,
			"config", Config,
		)
	} else {
		logger.Infow(
			context.Background(),
			"using default config",
			"config", Config,
		)
	}
}
