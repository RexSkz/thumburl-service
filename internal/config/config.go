package config

import (
	"thumburl-service/internal/pkg/cdpagent"
)

var PoolConfig = []*cdpagent.InitPoolConfig{
	{
		URL:          "http://localhost:9222",
		Count:        5,
		MaxUsedTimes: 10,
		TimeoutSec:   10,
	},
}
