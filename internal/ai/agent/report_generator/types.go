package report_generator

import (
	"SuperBizAgent/internal/model"
	"time"
)

// ReportRequest 报告生成请求
type ReportRequest struct {
	Type        model.ReportType `json:"type"`
	Title       string           `json:"title"`
	StartTime   time.Time        `json:"start_time"`
	EndTime     time.Time        `json:"end_time"`
	TemplateID  uint             `json:"template_id"`
	EventIDs    []uint           `json:"event_ids"`
	Keywords    []string         `json:"keywords"`
	GeneratedBy string           `json:"generated_by"`
}

// ReportResponse 报告生成响应
type ReportResponse struct {
	ReportID uint   `json:"report_id"`
	Title    string `json:"title"`
	Summary  string `json:"summary"`
	Content  string `json:"content"`
	Status   string `json:"status"`
	Error    string `json:"error,omitempty"`
}

// EventSummary 事件摘要
type EventSummary struct {
	TotalCount    int `json:"total_count"`
	CriticalCount int `json:"critical_count"`
	HighCount     int `json:"high_count"`
	MediumCount   int `json:"medium_count"`
	LowCount      int `json:"low_count"`
}
