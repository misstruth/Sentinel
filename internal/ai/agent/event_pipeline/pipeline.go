package event_pipeline

import (
	"context"
	"time"

	"SuperBizAgent/internal/ai/agent/event_dedup"
	"SuperBizAgent/internal/ai/agent/event_extraction"
	"SuperBizAgent/internal/ai/agent/risk_assessment"
)

// Pipeline 事件处理流水线
type Pipeline struct {
	ctx context.Context
}

func NewPipeline(ctx context.Context) *Pipeline {
	return &Pipeline{ctx: ctx}
}

// Process 处理事件流水线: 提取 -> 去重 -> 风险评估
func (p *Pipeline) Process(input *PipelineInput) (*PipelineResult, error) {
	var allEvents []*event_extraction.ExtractedEvent

	// Step 1: 提取
	extractor := event_extraction.NewExtractor(p.ctx)
	for _, raw := range input.RawEvents {
		result, _ := extractor.Extract(raw)
		if result.Success {
			allEvents = append(allEvents, result.Events...)
		}
	}

	// Step 2: 去重
	dedup := event_dedup.NewDeduplicator()
	dedupResult, _ := dedup.Dedup(p.ctx, &event_dedup.DedupInput{Events: allEvents})

	// Step 3: 风险评估
	assessor := risk_assessment.NewAssessor(p.ctx)
	var processed []*ProcessedEvent
	for _, ev := range dedupResult.UniqueEvents {
		assess, _ := assessor.Assess(&risk_assessment.AssessmentInput{Event: ev})
		processed = append(processed, &ProcessedEvent{
			ExtractedEvent: ev,
			RiskScore:      assess.RiskScore,
			Recommendation: assess.Recommendation,
		})
	}

	return &PipelineResult{
		Events:      processed,
		TotalCount:  len(allEvents),
		DedupCount:  dedupResult.DupCount,
		ProcessedAt: time.Now().Format(time.RFC3339),
	}, nil
}
