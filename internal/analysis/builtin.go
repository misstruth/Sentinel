package analysis

// 内置模板类别
const (
	CategoryVuln    = "vulnerability"
	CategoryThreat  = "threat"
	CategoryMalware = "malware"
)

const vulnPrompt = `分析以下CVE漏洞信息:
{{.Content}}

请提供:
1. 漏洞概述
2. 影响范围
3. 修复建议
`

const threatPrompt = `评估以下安全威胁:
{{.Content}}

请提供:
1. 威胁等级
2. 攻击向量
3. 防护措施
`

// InitBuiltinTemplates 初始化内置模板
func InitBuiltinTemplates(m *Manager) {
	m.Add(&Template{
		ID:          1,
		Name:        "漏洞分析",
		Description: "分析CVE漏洞详情",
		Prompt:      vulnPrompt,
		Category:    CategoryVuln,
		IsBuiltin:   true,
	})

	m.Add(&Template{
		ID:          2,
		Name:        "威胁评估",
		Description: "评估安全威胁等级",
		Prompt:      threatPrompt,
		Category:    CategoryThreat,
		IsBuiltin:   true,
	})
}
