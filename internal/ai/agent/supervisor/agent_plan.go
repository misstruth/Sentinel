package supervisor

import (
	"SuperBizAgent/internal/ai/agent/plan_execute_replan"
	"context"
	"strings"
)

type PlanAgent struct{}

func (a *PlanAgent) Name() AgentType { return AgentPlan }

func (a *PlanAgent) Execute(ctx context.Context, task *Task, cb StreamCallback) (*Result, error) {
	cb(AgentPlan, "[Plan Agent 规划执行...]\n")

	result, details, err := plan_execute_replan.BuildPlanAgent(ctx, task.Query)
	if err != nil {
		return &Result{TaskID: task.ID, Agent: AgentPlan, Error: err}, err
	}

	for _, d := range details {
		cb(AgentPlan, d+"\n")
	}

	content := strings.Join(details, "\n") + "\n" + result
	return &Result{TaskID: task.ID, Agent: AgentPlan, Content: content}, nil
}

func init() { RegisterAgent(&PlanAgent{}) }
