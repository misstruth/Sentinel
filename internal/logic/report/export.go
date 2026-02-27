package report

import (
	"SuperBizAgent/internal/model"
	"bytes"
	"encoding/json"
	"fmt"
)

// ExportService 导出服务
type ExportService struct{}

// NewExportService 创建导出服务
func NewExportService() *ExportService {
	return &ExportService{}
}

// ExportFormat 导出格式
type ExportFormat string

const (
	FormatMarkdown ExportFormat = "markdown"
	FormatHTML     ExportFormat = "html"
	FormatJSON     ExportFormat = "json"
)

// Export 导出报告
func (s *ExportService) Export(report *model.Report, format ExportFormat) ([]byte, error) {
	switch format {
	case FormatMarkdown:
		return s.exportMarkdown(report)
	case FormatHTML:
		return s.exportHTML(report)
	case FormatJSON:
		return s.exportJSON(report)
	default:
		return nil, fmt.Errorf("不支持的导出格式: %s", format)
	}
}

// exportMarkdown 导出为Markdown格式
func (s *ExportService) exportMarkdown(report *model.Report) ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("# %s\n\n", report.Title))
	buf.WriteString(fmt.Sprintf("**类型**: %s\n", report.Type))
	buf.WriteString(fmt.Sprintf("**生成时间**: %s\n\n", report.CreatedAt.Format("2006-01-02 15:04:05")))
	buf.WriteString("---\n\n")
	buf.WriteString(report.Content)

	return buf.Bytes(), nil
}

// exportJSON 导出为JSON格式
func (s *ExportService) exportJSON(report *model.Report) ([]byte, error) {
	data := map[string]interface{}{
		"id":         report.ID,
		"title":      report.Title,
		"type":       report.Type,
		"content":    report.Content,
		"summary":    report.Summary,
		"created_at": report.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return json.Marshal(data)
}

// exportHTML 导出为HTML格式
func (s *ExportService) exportHTML(report *model.Report) ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString("<!DOCTYPE html>\n<html>\n<head>\n")
	buf.WriteString("<meta charset=\"UTF-8\">\n")
	buf.WriteString(fmt.Sprintf("<title>%s</title>\n", report.Title))
	buf.WriteString("<style>body{font-family:sans-serif;max-width:800px;margin:0 auto;padding:20px;}</style>\n")
	buf.WriteString("</head>\n<body>\n")
	buf.WriteString(fmt.Sprintf("<h1>%s</h1>\n", report.Title))
	buf.WriteString(fmt.Sprintf("<p><strong>类型:</strong> %s</p>\n", report.Type))
	buf.WriteString(fmt.Sprintf("<p><strong>生成时间:</strong> %s</p>\n", report.CreatedAt.Format("2006-01-02 15:04:05")))
	buf.WriteString("<hr>\n")
	buf.WriteString(fmt.Sprintf("<div>%s</div>\n", report.Content))
	buf.WriteString("</body>\n</html>")

	return buf.Bytes(), nil
}
