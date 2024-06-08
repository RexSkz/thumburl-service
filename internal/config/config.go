package config

import (
	"thumburl-service/internal/pkg/cdpagent"
)

var PoolConfig = []*cdpagent.InitPoolConfig{
	{
		URL:          "http://localhost:9222",
		Count:        10,
		MaxUsedTimes: 20,
		TimeoutSec:   10,
	},
}
