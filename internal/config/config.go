package config

import (
	"time"
)

type InitPoolConfig struct {
	URL          string
	Count        int
	MaxUsedTimes int
	TimeoutSec   int
}

const Port = ":8080"

var PoolConfig = []*InitPoolConfig{
	{
		URL:          "http://localhost:9222",
		Count:        10,
		MaxUsedTimes: 20,
		TimeoutSec:   10,
	},
}

const GetAgentTimeout = 5 * time.Second
