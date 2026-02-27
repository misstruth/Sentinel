package report_generator

// DailyReportPrompt 日报生成提示词
const DailyReportPrompt = `你是一个专业的安全分析师，请根据以下安全事件数据生成一份日报。

## 报告要求
1. 使用 Markdown 格式
2. 包含执行摘要、重点事件、统计分析、建议措施
3. 语言简洁专业

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

## 事件数据
{{.EventData}}

## 时间范围
{{.StartTime}} 至 {{.EndTime}}

请生成报告：`
