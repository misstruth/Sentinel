package supervisor

import (
	"context"
	"encoding/json"
	"fmt"

	"SuperBizAgent/internal/ai/models"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

// RouterOutput 路由节点输出
type RouterOutput struct {
	Input     *SupervisorInput
	AgentType AgentType
}

// newRouterLambda 创建路由 Lambda 节点
func newRouterLambda(ctx context.Context) (*compose.Lambda, error) {
	model, err := models.OpenAIForDeepSeekV3Quick(ctx)
	if err != nil {
		return nil, fmt.Errorf("create router model: %w", err)
	}

	return compose.InvokableLambda(func(ctx context.Context, input *SupervisorInput) (*RouterOutput, error) {
		agentType, err := classifyIntent(ctx, model, input.Query)
		if err != nil {
			// 分类失败时降级到 chat
			agentType = AgentChat
		}

		// 回调通知路由结果
		if input.Callback != nil {
			input.Callback(agentType, fmt.Sprintf("[路由到: %s]\n", agentType))
		}

		return &RouterOutput{
			Input:     input,
			AgentType: agentType,
		}, nil
	}), nil
}

// classifyIntent 使用 LLM 分类用户意图
func classifyIntent(ctx context.Context, m model.ToolCallingChatModel, query string) (AgentType, error) {
	prompt := buildClassifyPrompt(query)

	resp, err := m.Generate(ctx, []*schema.Message{schema.UserMessage(prompt)})
	if err != nil {
		return AgentChat, fmt.Errorf("classify intent: %w", err)
	}

	return parseAgentType(resp.Content), nil
}

func buildClassifyPrompt(query string) string {
	return `你是路由分类器。根据用户问题选择最合适的Agent。

<agents>
- chat: 通用对话、安全咨询、知识问答、日志查询、订阅管理（默认选项）
- event: 安全事件查询、事件分析、事件关联、告警关联、事件处置建议
- report: 生成报告、查看报告、报告数据分析
- risk: 风险评估、威胁分析、漏洞评估、安全评分、CVE分析
- plan: 复杂多步骤任务、需要规划的操作
</agents>

<examples>
"最近有什么安全事件" → event
"有哪些高危漏洞" → event
"事件123的详情" → event
"这些事件之间有什么关联" → event
"告警和事件有什么关系" → event
"帮我生成本周安全报告" → report
"查看最近的报告" → report
"月报数据分析" → report
"评估这个CVE的风险" → risk
"当前系统的安全评分" → risk
"这个漏洞的威胁有多大" → risk
"帮我部署一套监控系统" → plan
"查一下日志" → chat
"怎么配置告警规则" → chat
</examples>

用户问题: ` + query + `

只返回JSON: {"agent": "xxx"}`
}

func parseAgentType(content string) AgentType {
	var result struct {
		Agent string `json:"agent"`
	}

	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return AgentChat
	}

	switch result.Agent {
	case "event":
		return AgentEvent
	case "report":
		return AgentReport
	case "risk":
		return AgentRisk
	case "plan":
		return AgentPlan
	default:
		return AgentChat
	}
}
