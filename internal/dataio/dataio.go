package dataio

import (
	"encoding/json"
	"io"
)

// Format 数据格式
type Format string

const (
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
)

// Exporter 导出器
type Exporter struct {
	format Format
}

// NewExporter 创建导出器
func NewExporter(f Format) *Exporter {
	return &Exporter{format: f}
}

// Export 导出数据
func (e *Exporter) Export(w io.Writer, data interface{}) error {
	if e.format == FormatJSON {
		return json.NewEncoder(w).Encode(data)
	}
	return nil
}

// Importer 导入器
type Importer struct {
	format Format
}

// NewImporter 创建导入器
func NewImporter(f Format) *Importer {
	return &Importer{format: f}
}

// Import 导入数据
func (i *Importer) Import(r io.Reader, v interface{}) error {
	if i.format == FormatJSON {
		return json.NewDecoder(r).Decode(v)
	}
	return nil
}
