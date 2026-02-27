package report_generator

// DailyReportPrompt 日报生成提示词
const DailyReportPrompt = `你是一个专业的安全分析师，请根据以下安全事件数据生成一份日报。

## 报告要求
1. 使用 Markdown 格式
2. 包含执行摘要、重点事件、统计分析、建议措施
3. 语言简洁专业

## 统计概览
- 事件总数: {{.TotalCount}}
- 严重: {{.CriticalCount}}, 高危: {{.HighCount}}, 中危: {{.MediumCount}}, 低危: {{.LowCount}}

## 事件数据
{{.EventData}}

## 时间范围
{{.StartTime}} 至 {{.EndTime}}

请生成报告：`

// WeeklyReportPrompt 周报生成提示词
const WeeklyReportPrompt = `你是一个专业的安全分析师，请根据以下安全事件数据生成一份周报。

## 报告要求
1. 使用 Markdown 格式
2. 包含本周概览、趋势分析、重点事件、下周关注点
3. 对比上周数据变化

## 统计概览
- 事件总数: {{.TotalCount}}
- 严重: {{.CriticalCount}}, 高危: {{.HighCount}}, 中危: {{.MediumCount}}, 低危: {{.LowCount}}

## 事件数据
{{.EventData}}

## 时间范围
{{.StartTime}} 至 {{.EndTime}}

请生成报告：`

// MonthlyReportPrompt 月报生成提示词
const MonthlyReportPrompt = `你是一个专业的安全分析师，请根据以下安全事件数据生成一份月度安全报告。

## 报告要求
1. 使用 Markdown 格式
2. 包含月度概览、趋势分析（按周对比）、重点漏洞/威胁、各严重等级分布、改进建议
3. 重点分析 CVE 漏洞的影响范围和修复优先级
4. 提供下月安全工作建议

## 统计概览
- 事件总数: {{.TotalCount}}
- 严重: {{.CriticalCount}}, 高危: {{.HighCount}}, 中危: {{.MediumCount}}, 低危: {{.LowCount}}

## 事件数据
{{.EventData}}

## 时间范围
{{.StartTime}} 至 {{.EndTime}}

请生成报告：`

// VulnAlertPrompt 漏洞告警生成提示词
const VulnAlertPrompt = `你是一个专业的漏洞分析师，请根据以下安全事件数据生成一份漏洞告警报告。

## 报告要求
1. 使用 Markdown 格式
2. 聚焦漏洞分析：CVE编号、CVSS评分、受影响厂商/产品、利用难度
3. 按 CVSS 评分从高到低排列
4. 每个漏洞给出修复建议和临时缓解措施
5. 包含漏洞影响评估和修复优先级矩阵

## 统计概览
- 事件总数: {{.TotalCount}}
- 严重: {{.CriticalCount}}, 高危: {{.HighCount}}, 中危: {{.MediumCount}}, 低危: {{.LowCount}}

## 事件数据
{{.EventData}}

## 时间范围
{{.StartTime}} 至 {{.EndTime}}

请生成报告：`

// ThreatBriefPrompt 威胁简报生成提示词
const ThreatBriefPrompt = `你是一个专业的威胁情报分析师，请根据以下安全事件数据生成一份威胁简报。

## 报告要求
1. 使用 Markdown 格式
2. 包含威胁态势概览、活跃威胁组织/攻击手法、IOC指标汇总
3. 分析攻击趋势和目标行业
4. 提供防御建议和检测规则建议
5. 语言简洁，适合快速阅读

## 统计概览
- 事件总数: {{.TotalCount}}
- 严重: {{.CriticalCount}}, 高危: {{.HighCount}}, 中危: {{.MediumCount}}, 低危: {{.LowCount}}

## 事件数据
{{.EventData}}

## 时间范围
{{.StartTime}} 至 {{.EndTime}}

请生成报告：`

// CustomReportPrompt 自定义/决策简报生成提示词
const CustomReportPrompt = `你是一个专业的安全分析师，请根据以下安全事件数据生成一份安全分析报告。

## 报告要求
1. 使用 Markdown 格式
2. 包含执行摘要（一段话概括核心风险）、事件分析、风险评估、处置建议
3. 适合管理层快速阅读，重点突出，结论明确
4. 如有 CVE 漏洞，标注修复优先级

## 统计概览
- 事件总数: {{.TotalCount}}
- 严重: {{.CriticalCount}}, 高危: {{.HighCount}}, 中危: {{.MediumCount}}, 低危: {{.LowCount}}

## 事件数据
{{.EventData}}

## 时间范围
{{.StartTime}} 至 {{.EndTime}}

请生成报告：`
