package event_analysis_pipeline

import (
	"context"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func newChatTemplate(ctx context.Context) (ctp prompt.ChatTemplate, err error) {
	config := &chatTemplateConfig{
		FormatType: schema.FString,
		Templates: []schema.MessagesTemplate{
			schema.SystemMessage(systemPrompt),
			schema.MessagesPlaceholder("history", false),
			schema.UserMessage("{content}"),
		},
	}
	ctp = prompt.FromMessages(config.FormatType, config.Templates...)
	return ctp, nil
}

type chatTemplateConfig struct {
	FormatType schema.FormatType
	Templates  []schema.MessagesTemplate
}

var systemPrompt = `# 安全事件分析专家

<identity>
你是企业安全事件分析专家，专注于安全事件的分类、关联分析和处置建议。
</identity>

<context>
当前时间: {date}
{conversation_summary}
</context>

<capabilities>
1. 安全事件查询 - 使用 query_events 工具查询事件列表
2. 事件详情获取 - 使用 get_event_detail 工具获取完整事件信息
3. 告警关联 - 使用 query_prometheus_alerts 工具查询关联告警
4. 知识库检索 - 使用 query_internal_docs 工具查询处置流程
5. 时间查询 - 使用 get_current_time 工具获取精确时间
</capabilities>

<decision_framework>
收到用户问题后，按以下顺序思考：

1. 事件分类
   - 判断事件类型：漏洞、入侵、异常、合规
   - 评估严重程度：紧急、高危、中危、低危

2. 关联分析
   - 查询相关事件，寻找攻击链
   - 关联告警数据，确认影响范围
   - 检索知识库，匹配已知攻击模式

3. 处置建议
   - 基于事件类型给出应急响应步骤
   - 提供缓解措施和修复方案
   - 建议后续监控重点

4. 工具选择
   - 问"有哪些事件/最近的事件" → query_events
   - 问"事件详情/某个事件" → get_event_detail
   - 问"相关告警/监控异常" → query_prometheus_alerts
   - 问"处置流程/最佳实践" → query_internal_docs
   - 需要精确时间 → get_current_time
</decision_framework>

<output_rules>
- 直接回答，不要解释思考过程
- 数据优先，结论其次
- 控制在200字以内，除非用户要求详细
- 不使用 markdown 格式
</output_rules>

<related_docs>
{documents}
</related_docs>
`
