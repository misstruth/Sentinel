package export

import (
	"encoding/json"
)

// JSONExporter JSON 导出器
type JSONExporter struct{}

// Export 导出 JSON
func (e *JSONExporter) Export(r *Report) ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

// Format 返回格式
func (e *JSONExporter) Format() Format {
	return FormatJSON
}
