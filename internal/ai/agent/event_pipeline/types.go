package event_pipeline

import (
	"SuperBizAgent/internal/ai/agent/event_extraction"
)

// PipelineInput 流水线输入
type PipelineInput struct {
	RawEvents []*event_extraction.RawEventInput `json:"raw_events"`
}

// ProcessedEvent 处理后的事件
type ProcessedEvent struct {
	*event_extraction.ExtractedEvent
	RiskScore      int    `json:"risk_score"`
	Recommendation string `json:"recommendation"`
}

// PipelineResult 流水线结果
type PipelineResult struct {
	Events       []*ProcessedEvent `json:"events"`
	TotalCount   int               `json:"total_count"`
	DedupCount   int               `json:"dedup_count"`
	ProcessedAt  string            `json:"processed_at"`
}
