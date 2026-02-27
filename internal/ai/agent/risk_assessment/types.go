package risk_assessment

import "SuperBizAgent/internal/ai/agent/event_extraction"

// AssessmentInput 评估输入
type AssessmentInput struct {
	Event *event_extraction.ExtractedEvent `json:"event"`
	Query string                           `json:"query"`
}

// AssessmentResult 评估结果
type AssessmentResult struct {
	RiskScore      int      `json:"risk_score"`      // 0-100
	Severity       string   `json:"severity"`        // critical/high/medium/low/info
	Recommendation string   `json:"recommendation"`  // 处置建议
	Factors        []string `json:"factors"`         // 评分因素
}
