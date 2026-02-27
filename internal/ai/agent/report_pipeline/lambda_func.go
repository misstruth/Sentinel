package report_pipeline

import (
	"context"
	"time"
)

func newInputToRagLambda(ctx context.Context, input *UserMessage, opts ...any) (output string, err error) {
	return input.Query, nil
}

func newInputToChatLambda(ctx context.Context, input *UserMessage, opts ...any) (output map[string]any, err error) {
	return map[string]any{
		"content":              input.Query,
		"history":              input.History,
		"date":                 time.Now().Format("2006-01-02 15:04:05"),
		"conversation_summary": "",
	}, nil
}
