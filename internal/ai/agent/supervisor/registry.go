package supervisor

import (
	"context"
	"sync"
)

// SubAgent 子 Agent 接口
type SubAgent interface {
	Name() AgentType
	Execute(ctx context.Context, task *Task, cb StreamCallback) (*Result, error)
}

var (
	agents = make(map[AgentType]SubAgent)
	agentMu sync.RWMutex
)

func RegisterAgent(agent SubAgent) {
	agentMu.Lock()
	defer agentMu.Unlock()
	agents[agent.Name()] = agent
}

func GetAgent(t AgentType) SubAgent {
	agentMu.RLock()
	defer agentMu.RUnlock()
	return agents[t]
}

func ListAgents() []AgentType {
	agentMu.RLock()
	defer agentMu.RUnlock()
	list := make([]AgentType, 0, len(agents))
	for k := range agents {
		list = append(list, k)
	}
	return list
}
