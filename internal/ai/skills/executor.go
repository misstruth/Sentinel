package skills

import (
	"SuperBizAgent/internal/ai/models"
	"SuperBizAgent/internal/ai/tools"
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

// Executor Skill 执行器
type Executor struct {
	skill  *Skill
	params map[string]any
}

// NewExecutor 创建执行器
func NewExecutor(skillID string, params map[string]any) (*Executor, error) {
	skill := Get(skillID)
	if skill == nil {
		return nil, fmt.Errorf("skill not found: %s", skillID)
	}
	return &Executor{skill: skill, params: params}, nil
}

// Execute 执行 Skill
func (e *Executor) Execute(ctx context.Context, callback func(ExecuteResult)) error {
	callback(ExecuteResult{Type: "step", Content: "正在准备执行 " + e.skill.Name + "..."})

	prompt := e.buildPrompt()
	selectedTools := e.selectTools()

	model, err := newModel(ctx)
	if err != nil {
		return err
	}

	config := &react.AgentConfig{
		MaxStep:          10,
		ToolCallingModel: model,
	}
	config.ToolsConfig.Tools = selectedTools

	agent, err := react.NewAgent(ctx, config)
	if err != nil {
		return err
	}

	callback(ExecuteResult{Type: "step", Content: "正在分析..."})

	result, err := agent.Generate(ctx, []*schema.Message{schema.UserMessage(prompt)})
	if err != nil {
		return err
	}

	callback(ExecuteResult{Type: "result", Content: result.Content})
	return nil
}

func (e *Executor) buildPrompt() string {
	prompt := e.skill.Prompt
	for k, v := range e.params {
		prompt = strings.ReplaceAll(prompt, "{"+k+"}", fmt.Sprintf("%v", v))
	}
	return prompt
}

func (e *Executor) selectTools() []tool.BaseTool {
	toolMap := map[string]tool.BaseTool{
		"query_events":        tools.NewEventQueryTool(),
		"query_subscriptions": tools.NewSubscriptionQueryTool(),
		"query_reports":       tools.NewReportQueryTool(),
		"query_internal_docs": tools.NewQueryInternalDocsTool(),
		"get_current_time":    tools.NewGetCurrentTimeTool(),
		"mysql_crud":          tools.NewMysqlCrudTool(),
	}

	var selected []tool.BaseTool
	for _, name := range e.skill.Tools {
		if t, ok := toolMap[name]; ok {
			selected = append(selected, t)
		}
	}
	return selected
}

func newModel(ctx context.Context) (model.ToolCallingChatModel, error) {
	return models.OpenAIForDeepSeekV3Quick(ctx)
}
