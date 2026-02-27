package supervisor

import (
	"context"
	"encoding/json"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// AgentTool 让 Agent 可作为 Tool 被调用
type AgentTool struct {
	agent SubAgent
}

func NewAgentTool(agentType AgentType) tool.BaseTool {
	return &AgentTool{agent: GetAgent(agentType)}
}

func (t *AgentTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "call_" + string(t.agent.Name()),
		Desc: "调用" + string(t.agent.Name()) + "Agent处理任务",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"query": {Type: schema.String, Desc: "任务描述", Required: true},
		}),
	}, nil
}

func (t *AgentTool) InvokableRun(ctx context.Context, args string, opts ...tool.Option) (string, error) {
	var p struct{ Query string }
	json.Unmarshal([]byte(args), &p)

	task := &Task{ID: "sub", Query: p.Query, Type: t.agent.Name()}
	result, err := t.agent.Execute(ctx, task, func(a AgentType, c string) {})
	if err != nil {
		return "", err
	}
	return result.Content, nil
}
