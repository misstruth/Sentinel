package template

// 预定义模板名称
const (
	AlertTemplate  = "alert"
	ReportTemplate = "report"
	DigestTemplate = "digest"
)

const alertBody = `
<h2>安全告警通知</h2>
<p><strong>标题:</strong> {{.Title}}</p>
<p><strong>级别:</strong> {{.Severity}}</p>
<p><strong>来源:</strong> {{.Source}}</p>
<p><strong>时间:</strong> {{.Time}}</p>
<hr>
<p>{{.Description}}</p>
`

const reportBody = `
<h2>安全分析报告</h2>
<p><strong>报告标题:</strong> {{.Title}}</p>
<p><strong>生成时间:</strong> {{.Time}}</p>
<hr>
<div>{{.Content}}</div>
`

const digestBody = `
<h2>每日安全摘要</h2>
<p><strong>日期:</strong> {{.Date}}</p>
<p><strong>事件总数:</strong> {{.Total}}</p>
<p><strong>高危事件:</strong> {{.Critical}}</p>
<hr>
<h3>事件列表</h3>
{{.EventList}}
`

// InitDefaultTemplates 初始化默认模板
func InitDefaultTemplates(e *Engine) {
	e.Register(&Template{
		Name:    AlertTemplate,
		Subject: "[安全告警] {{.Title}}",
		Body:    alertBody,
	})

	e.Register(&Template{
		Name:    ReportTemplate,
		Subject: "[安全报告] {{.Title}}",
		Body:    reportBody,
	})

	e.Register(&Template{
		Name:    DigestTemplate,
		Subject: "[每日摘要] {{.Date}} 安全事件汇总",
		Body:    digestBody,
	})
}
