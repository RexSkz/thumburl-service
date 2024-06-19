package cdpagent

import (
	"context"
	"strconv"
	"sync"
	"time"

	"thumburl-service/internal/config"
	"thumburl-service/internal/pkg/logger"

	"github.com/pkg/errors"
)

type agentKey string

type PoolAgentInfo struct {
	Agent        *Agent
	key          agentKey
	url          string
	currentInUse bool
	usedTimes    int
	maxUsedTimes int
	timeoutSec   int
}

type Pool struct {
	mu        sync.RWMutex
	agents    map[agentKey]*PoolAgentInfo
	available chan agentKey
}

func InitPool(configs []*config.PoolConfig) (*Pool, error) {
	var pool = new(Pool)

	totalAgents := 0
	for _, config := range configs {
		totalAgents += config.Count
	}

	pool.mu.Lock()
	defer pool.mu.Unlock()

	pool.agents = make(map[agentKey]*PoolAgentInfo)
	pool.available = make(chan agentKey, totalAgents)

	for _, config := range configs {
		for i := 0; i < config.Count; i++ {
			agent, err := newAgent(config.URL, config.TimeoutSec)
			if err != nil {
				return nil, errors.Wrap(err, "failed to create agent")
			}
			key := agentKey(config.URL + "/" + strconv.Itoa(i))
			pool.agents[key] = &PoolAgentInfo{
				Agent:        agent,
				key:          key,
				url:          config.URL,
				currentInUse: false,
				usedTimes:    0,
				maxUsedTimes: config.MaxUsedTimes,
				timeoutSec:   config.TimeoutSec,
			}
			pool.available <- key
			logger.Infow(
				context.Background(),
				"created agent",
				"index", i,
				"url", config.URL,
				"key", key,
			)
		}
	}

	return pool, nil
}

func (pool *Pool) GetAgent(ctx context.Context) (*PoolAgentInfo, error) {
	timeoutSec := time.Second * time.Duration(config.Config.GetAgentTimeoutSec)
	select {
	case key := <-pool.available:
		pool.mu.RLock()
		agent := pool.agents[key]
		pool.mu.RUnlock()

		if agent == nil {
			return nil, errors.New("agent '" + string(key) + "' not found")
		} else {
			logger.Infow(
				ctx,
				"agent retrieved",
				"key", key,
				"remain", len(pool.available),
			)
		}
		agent.currentInUse = true
		agent.usedTimes++
		return agent, nil
	case <-time.After(timeoutSec):
		return nil, errors.New("get agent timeout, no agent is released")
	}
}

func (pool *Pool) ReleaseAgent(ctx context.Context, agent *PoolAgentInfo, force bool) error {
	if !agent.currentInUse {
		return errors.New("agent '" + string(agent.key) + "' is not in use")
	}
	agent.currentInUse = false

	if agent.usedTimes > agent.maxUsedTimes || force {
		nextAgent, err := newAgent(agent.url, agent.timeoutSec)
		if err != nil {
			return errors.Wrap(err, "release agent failed")
		}

		pool.agents[agent.key].Agent.close()

		pool.mu.Lock()
		defer pool.mu.Unlock()
		pool.agents[agent.key] = &PoolAgentInfo{
			Agent:        nextAgent,
			key:          agent.key,
			url:          agent.url,
			currentInUse: false,
			usedTimes:    0,
			maxUsedTimes: agent.maxUsedTimes,
			timeoutSec:   agent.timeoutSec,
		}
	}

	pool.available <- agent.key
	if force {
		logger.Infow(
			ctx,
			"agent force released",
			"key", agent.key,
		)
	} else {
		logger.Infow(
			ctx,
			"agent released",
			"key", agent.key,
		)
	}
	return nil
}

func (pool *Pool) Dispose() {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	for _, agent := range pool.agents {
		if err := agent.Agent.close(); err != nil {
			logger.Errorw(
				context.Background(),
				"dispose agent failed",
				"key", agent.key,
				"error", err,
			)
		} else {
			logger.Infow(
				context.Background(),
				"disposed agent",
				"key", agent.key,
			)
		}
	}
}
