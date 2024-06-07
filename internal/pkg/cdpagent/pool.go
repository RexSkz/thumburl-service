package cdpagent

import (
	"fmt"
	"strconv"
	"sync"
	"time"

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

type InitPoolConfig struct {
	URL          string
	Count        int
	MaxUsedTimes int
	TimeoutSec   int
}

func InitPool(configs []*InitPoolConfig) (*Pool, error) {
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
			fmt.Printf("create agent #%d of %s\n", i, config.URL)
			agent, err := newAgent(config.URL, config.TimeoutSec)
			if err != nil {
				return nil, err
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
		}
	}

	return pool, nil
}

func (pool *Pool) GetAgent() (*PoolAgentInfo, error) {
	select {
	case key := <-pool.available:
		pool.mu.RLock()
		agent := pool.agents[key]
		pool.mu.RUnlock()

		if agent == nil {
			return nil, errors.New("agent not found")
		} else {
			fmt.Printf("get agent %s\n", key)
		}
		agent.currentInUse = true
		agent.usedTimes++
		return agent, nil
	case <-time.After(5 * time.Second):
		return nil, errors.New("timeout to get agent")
	}
}

func (pool *Pool) ReleaseAgent(agent *PoolAgentInfo) error {
	if !agent.currentInUse {
		return errors.New("agent is not in use")
	}
	agent.currentInUse = false

	if agent.usedTimes > agent.maxUsedTimes {
		nextAgent, err := newAgent(agent.url, agent.timeoutSec)
		if err != nil {
			return errors.Wrap(err, "failed to release agent")
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
	fmt.Printf("release agent %s\n", agent.key)
	return nil
}

func (pool *Pool) Dispose() {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	for _, agent := range pool.agents {
		fmt.Printf("dispose agent %s\n", agent.key)
		agent.Agent.close()
	}
}
