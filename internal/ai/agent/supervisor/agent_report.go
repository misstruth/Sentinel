package supervisor

import (
	"SuperBizAgent/internal/ai/agent/report_pipeline"
	"context"
)

type ReportAgent struct{}

func (a *ReportAgent) Name() AgentType { return AgentReport }

func (a *ReportAgent) Execute(ctx context.Context, task *Task, cb StreamCallback) (*Result, error) {
	cb(AgentReport, "[Report Agent 处理中...]\n")

	runner, err := report_pipeline.BuildReportAgent(ctx)
	if err != nil {
		return &Result{TaskID: task.ID, Agent: AgentReport, Error: err}, err
	}

	mem := GetSharedMemory()
	input := &report_pipeline.UserMessage{
		ID:      task.ID,
		Query:   task.Query,
		History: mem.GetMessages(),
	}

	out, err := runner.Invoke(ctx, input)
	if err != nil {
		return &Result{TaskID: task.ID, Agent: AgentReport, Error: err}, err
	}

	cb(AgentReport, out.Content)
	return &Result{TaskID: task.ID, Agent: AgentReport, Content: out.Content}, nil
}

func init() { RegisterAgent(&ReportAgent{}) }
