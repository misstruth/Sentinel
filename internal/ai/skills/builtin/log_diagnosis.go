package builtin

import "SuperBizAgent/internal/ai/skills"

var LogDiagnosis = &skills.Skill{
	ID:          "log_diagnosis",
	Name:        "日志诊断",
	Description: "分析日志数据，发现异常模式和潜在问题",
	Category:    "ops",
	Enabled:     true,
	Tools:       []string{"query_internal_docs", "get_current_time"},
	Params: []skills.SkillParam{
		{Name: "keyword", Type: "string", Description: "搜索关键词", Required: true},
	},
	Prompt: `你是运维专家。请根据关键词 "{keyword}" 进行日志诊断：
1. 分析可能的异常模式
2. 识别潜在问题
3. 提供排查建议

输出纯文本，不使用markdown。`,
}
