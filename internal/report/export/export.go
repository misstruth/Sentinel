package export

// Format 导出格式
type Format string

const (
	FormatPDF      Format = "pdf"
	FormatHTML     Format = "html"
	FormatMarkdown Format = "markdown"
	FormatJSON     Format = "json"
)

// Report 报告数据
type Report struct {
	Title   string
	Content string
	Time    string
}

// Exporter 导出器接口
type Exporter interface {
	Export(r *Report) ([]byte, error)
	Format() Format
}
