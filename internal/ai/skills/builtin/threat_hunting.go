package builtin

import "SuperBizAgent/internal/ai/skills"

var ThreatHunting = &skills.Skill{
	ID:          "threat_hunting",
	Name:        "威胁狩猎",
	Description: "主动搜索潜在威胁和攻击迹象",
	Category:    "security",
	Enabled:     true,
	Tools:       []string{"query_events", "query_subscriptions"},
	Params: []skills.SkillParam{
		{Name: "target", Type: "string", Description: "狩猎目标(IP/域名/用户)", Required: true},
	},
	Prompt: `你是威胁狩猎专家。请针对目标 "{target}" 进行威胁狩猎：
1. 查询相关安全事件
2. 分析可疑行为模式
3. 评估威胁等级

输出纯文本，不使用markdown。`,
}
