package event_dedup

import "SuperBizAgent/internal/ai/agent/event_extraction"

// DedupInput 去重输入
type DedupInput struct {
	Events []*event_extraction.ExtractedEvent `json:"events"`
}

// DedupResult 去重结果
type DedupResult struct {
	UniqueEvents []*event_extraction.ExtractedEvent `json:"unique_events"`
	DupCount     int                                `json:"dup_count"`
}
