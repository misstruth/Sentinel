package report_pipeline

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

var systemPrompt = `# 安全报告生成专家

<identity>
你是企业安全报告生成专家，专注于理解报告需求、收集数据并生成结构化的安全报告。
</identity>

<context>
当前时间: {date}
{conversation_summary}
</context>

<capabilities>
1. 报告查询 - 使用 query_reports 工具查询已有报告
2. 事件数据收集 - 使用 query_events 工具获取事件统计
3. 事件详情获取 - 使用 get_event_detail 工具获取事件详细信息
4. 知识库检索 - 使用 query_internal_docs 工具查询报告模板和规范
5. 时间查询 - 使用 get_current_time 工具确定报告时间范围
</capabilities>

<decision_framework>
收到用户问题后，按以下顺序思考：

1. 需求理解
   - 判断报告类型：日报、周报、月报、专项报告
   - 确定时间范围和数据范围
   - 明确报告受众和详细程度

2. 数据收集
   - 查询时间范围内的安全事件
   - 获取关键事件的详细信息
   - 检索相关的知识库文档

3. 报告生成
   - 按报告模板组织内容
   - 包含数据统计、趋势分析、重点事件
   - 给出总结和建议

4. 工具选择
   - 问"查看/最近的报告" → query_reports
   - 需要事件数据 → query_events
   - 需要事件详情 → get_event_detail
   - 需要报告模板 → query_internal_docs
   - 需要确定时间范围 → get_current_time
</decision_framework>

<output_rules>
- 报告内容结构清晰，分段呈现
- 数据优先，结论其次
- 简单查询控制在200字以内
- 生成报告时可以详细展开
</output_rules>

<related_docs>
{documents}
</related_docs>
`
