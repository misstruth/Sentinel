package supervisor

import (
	"context"
	"sync/atomic"

	"github.com/cloudwego/eino/schema"
)

var taskIDCounter atomic.Int64

// Supervisor 调度器
type Supervisor struct {
	ctx    context.Context
	memory *SharedMemory
}

func NewSupervisor(ctx context.Context) *Supervisor {
	return &Supervisor{ctx: ctx, memory: GetSharedMemory()}
}

// Execute 执行任务 - 使用 Graph 编排
func (s *Supervisor) Execute(query string, cb StreamCallback) (*Result, error) {
	// 记录用户消息
	s.memory.AddMessage(schema.UserMessage(query))

	// 构建并执行 Graph
	graph, err := BuildSupervisorGraph(s.ctx)
	if err != nil {
		return &Result{Error: err}, err
	}

	input := &SupervisorInput{
		Query:    query,
		Callback: cb,
	}

	output, err := graph.Invoke(s.ctx, input)
	if err != nil {
		return &Result{Error: err}, err
	}

	// 记录结果到 memory
	if output.Content != "" {
		s.memory.AddMessage(schema.AssistantMessage(output.Content, nil))
	}

	return &Result{
		TaskID:  output.TaskID,
		Agent:   output.Agent,
		Content: output.Content,
		Error:   output.Error,
	}, output.Error
}
