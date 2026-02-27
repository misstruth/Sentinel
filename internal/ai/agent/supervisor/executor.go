package supervisor

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/compose"
)

// newExecutorLambda 创建执行 Lambda 节点
func newExecutorLambda() *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, input *RouterOutput) (*SupervisorOutput, error) {
		agent := GetAgent(input.AgentType)
		if agent == nil {
			agent = GetAgent(AgentChat)
		}

		task := &Task{
			ID:    generateTaskID(),
			Query: input.Input.Query,
			Type:  input.AgentType,
		}

		result, err := agent.Execute(ctx, task, input.Input.Callback)
		if err != nil {
			return &SupervisorOutput{
				TaskID:  task.ID,
				Agent:   input.AgentType,
				Content: "",
				Error:   err,
			}, nil // 返回 nil error，让 Graph 继续执行
		}

		return &SupervisorOutput{
			TaskID:  result.TaskID,
			Agent:   result.Agent,
			Content: result.Content,
			Error:   result.Error,
		}, nil
	})
}

// generateTaskID 生成任务 ID
func generateTaskID() string {
	return fmt.Sprintf("task_%d", taskIDCounter.Add(1))
}
