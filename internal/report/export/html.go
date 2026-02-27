package export

import (
	"fmt"
)

// HTMLExporter HTML 导出器
type HTMLExporter struct{}

// Export 导出 HTML
func (e *HTMLExporter) Export(r *Report) ([]byte, error) {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><title>%s</title></head>
<body>
<h1>%s</h1>
<p>%s</p>
<div>%s</div>
</body>
</html>`, r.Title, r.Title, r.Time, r.Content)
	return []byte(html), nil
}

// Format 返回格式
func (e *HTMLExporter) Format() Format {
	return FormatHTML
}
