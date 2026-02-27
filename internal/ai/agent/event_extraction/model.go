package event_extraction

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
)

func newExtractionModel(ctx context.Context) (model.ChatModel, error) {
	return openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:  "bc499880-ede3-4023-8991-2e84c0a83dd1",
		Model:   "deepseek-v3-250324",
		BaseURL: "https://ark.cn-beijing.volces.com/api/v3",
	})
}
