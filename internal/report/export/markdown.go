package export

import (
	"fmt"
)

// MarkdownExporter Markdown 导出器
type MarkdownExporter struct{}

// Export 导出 Markdown
func (e *MarkdownExporter) Export(r *Report) ([]byte, error) {
	md := fmt.Sprintf("# %s\n\n**时间:** %s\n\n%s", r.Title, r.Time, r.Content)
	return []byte(md), nil
}

// Format 返回格式
func (e *MarkdownExporter) Format() Format {
	return FormatMarkdown
}
