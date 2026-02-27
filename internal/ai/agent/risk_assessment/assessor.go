package risk_assessment

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/cloudwego/eino/schema"
)

type Assessor struct {
	ctx context.Context
}

func NewAssessor(ctx context.Context) *Assessor {
	return &Assessor{ctx: ctx}
}

func (a *Assessor) Assess(input *AssessmentInput) (*AssessmentResult, error) {
	prompt, _ := a.renderPrompt(input)
	model, err := newAssessmentModel(a.ctx)
	if err != nil {
		return a.defaultResult(), nil
	}

	resp, err := model.Generate(a.ctx, []*schema.Message{schema.UserMessage(prompt)})
	if err != nil {
		return a.defaultResult(), nil
	}

	return a.parseResponse(resp.Content)
}

func (a *Assessor) renderPrompt(input *AssessmentInput) (string, error) {
	t, _ := template.New("assess").Parse(AssessmentPrompt)
	var buf bytes.Buffer
	t.Execute(&buf, map[string]string{
		"Title":            input.Event.Title,
		"Description":      input.Event.Description,
		"EventType":        input.Event.EventType,
		"CVEIDs":           strings.Join(input.Event.CVEIDs, ", "),
		"AffectedProducts": strings.Join(input.Event.AffectedProducts, ", "),
	})
	return buf.String(), nil
}

func (a *Assessor) parseResponse(content string) (*AssessmentResult, error) {
	var result AssessmentResult
	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")
	if start >= 0 && end > start {
		json.Unmarshal([]byte(content[start:end+1]), &result)
	}
	if result.RiskScore == 0 {
		return a.defaultResult(), nil
	}
	return &result, nil
}

func (a *Assessor) defaultResult() *AssessmentResult {
	return &AssessmentResult{RiskScore: 50, Severity: "medium", Recommendation: "需要进一步分析"}
}
