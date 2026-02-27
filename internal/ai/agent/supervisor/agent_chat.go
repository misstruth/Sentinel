package supervisor

import (
	"SuperBizAgent/internal/ai/agent/chat_pipeline"
	"context"
)

type ChatAgent struct{}

func (a *ChatAgent) Name() AgentType { return AgentChat }

func (a *ChatAgent) Execute(ctx context.Context, task *Task, cb StreamCallback) (*Result, error) {
	cb(AgentChat, "[Chat Agent 处理中...]\n")

	runner, err := chat_pipeline.BuildChatAgent(ctx)
	if err != nil {
		return &Result{TaskID: task.ID, Agent: AgentChat, Error: err}, err
	}

	mem := GetSharedMemory()
	input := &chat_pipeline.UserMessage{
		ID:      task.ID,
		Query:   task.Query,
		History: mem.GetMessages(),
	}

	out, err := runner.Invoke(ctx, input)
	if err != nil {
		return &Result{TaskID: task.ID, Agent: AgentChat, Error: err}, err
	}

	cb(AgentChat, out.Content)
	return &Result{TaskID: task.ID, Agent: AgentChat, Content: out.Content}, nil
}

func init() { RegisterAgent(&ChatAgent{}) }
