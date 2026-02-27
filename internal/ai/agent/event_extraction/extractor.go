package event_extraction

import (
	"bytes"
	"context"
	"encoding/json"
	"text/template"

	"github.com/cloudwego/eino/schema"
)

// Extractor 事件提取器
type Extractor struct {
	ctx context.Context
}

func NewExtractor(ctx context.Context) *Extractor {
	return &Extractor{ctx: ctx}
}

func (e *Extractor) Extract(input *RawEventInput) (*ExtractionResult, error) {
	prompt, err := e.renderPrompt(input)
	if err != nil {
		return &ExtractionResult{Success: false, Error: err.Error()}, nil
	}

	model, err := newExtractionModel(e.ctx)
	if err != nil {
		return &ExtractionResult{Success: false, Error: err.Error()}, nil
	}

	resp, err := model.Generate(e.ctx, []*schema.Message{
		schema.UserMessage(prompt),
	})
	if err != nil {
		return &ExtractionResult{Success: false, Error: err.Error()}, nil
	}

	return e.parseResponse(resp.Content, input)
}

func (e *Extractor) renderPrompt(input *RawEventInput) (string, error) {
	t, err := template.New("extraction").Parse(ExtractionPrompt)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, map[string]string{
		"RawContent": input.RawContent,
		"Source":     input.Source,
		"SourceURL":  input.SourceURL,
	})
	return buf.String(), err
}

func (e *Extractor) parseResponse(content string, input *RawEventInput) (*ExtractionResult, error) {
	var result ExtractionResult
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		// 尝试提取JSON部分
		start := bytes.Index([]byte(content), []byte("{"))
		end := bytes.LastIndex([]byte(content), []byte("}"))
		if start >= 0 && end > start {
			json.Unmarshal([]byte(content[start:end+1]), &result)
		}
	}
	// 补充来源信息
	for _, ev := range result.Events {
		ev.Source = input.Source
		ev.SourceURL = input.SourceURL
	}
	result.Success = true
	return &result, nil
}
