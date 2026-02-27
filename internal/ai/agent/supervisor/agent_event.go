package supervisor

import (
	"SuperBizAgent/internal/ai/agent/event_analysis_pipeline"
	"context"
)

type EventAgent struct{}

func (a *EventAgent) Name() AgentType { return AgentEvent }

func (a *EventAgent) Execute(ctx context.Context, task *Task, cb StreamCallback) (*Result, error) {
	cb(AgentEvent, "[Event Agent 分析中...]\n")

	runner, err := event_analysis_pipeline.BuildEventAnalysisAgent(ctx)
	if err != nil {
		return &Result{TaskID: task.ID, Agent: AgentEvent, Error: err}, err
	}

	mem := GetSharedMemory()
	input := &event_analysis_pipeline.UserMessage{
		ID:      task.ID,
		Query:   task.Query,
		History: mem.GetMessages(),
	}

	out, err := runner.Invoke(ctx, input)
	if err != nil {
		return &Result{TaskID: task.ID, Agent: AgentEvent, Error: err}, err
	}

	cb(AgentEvent, out.Content)
	return &Result{TaskID: task.ID, Agent: AgentEvent, Content: out.Content}, nil
}

func init() { RegisterAgent(&EventAgent{}) }
