package supervisor

import (
	"SuperBizAgent/internal/ai/agent/risk_pipeline"
	"context"
)

type RiskAgent struct{}

func (a *RiskAgent) Name() AgentType { return AgentRisk }

func (a *RiskAgent) Execute(ctx context.Context, task *Task, cb StreamCallback) (*Result, error) {
	cb(AgentRisk, "[Risk Agent 评估中...]\n")

	runner, err := risk_pipeline.BuildRiskAgent(ctx)
	if err != nil {
		return &Result{TaskID: task.ID, Agent: AgentRisk, Error: err}, err
	}

	mem := GetSharedMemory()
	input := &risk_pipeline.UserMessage{
		ID:      task.ID,
		Query:   task.Query,
		History: mem.GetMessages(),
	}

	out, err := runner.Invoke(ctx, input)
	if err != nil {
		return &Result{TaskID: task.ID, Agent: AgentRisk, Error: err}, err
	}

	cb(AgentRisk, out.Content)
	return &Result{TaskID: task.ID, Agent: AgentRisk, Content: out.Content}, nil
}

func init() { RegisterAgent(&RiskAgent{}) }
