package builtin

import "SuperBizAgent/internal/ai/skills"

func init() {
	skills.Register(EventAnalysis)
	skills.Register(LogDiagnosis)
	skills.Register(ThreatHunting)
}

var EventAnalysis = &skills.Skill{
	ID:          "event_analysis",
	Name:        "安全事件分析",
	Description: "对安全事件进行深度分析，包括威胁评估、影响范围、处置建议",
	Category:    "security",
	Enabled:     true,
	Tools:       []string{"query_events", "query_internal_docs"},
	Params: []skills.SkillParam{
		{Name: "event_id", Type: "number", Description: "事件ID", Required: true},
	},
	Prompt: `你是安全分析专家。请对事件ID {event_id} 进行深度分析：
1. 先查询该事件的详细信息
2. 分析威胁等级和影响范围
3. 提供处置建议

输出纯文本，不使用markdown。`,
}
