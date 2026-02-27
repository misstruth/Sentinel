package report_generator

import (
	"SuperBizAgent/internal/ai/models"
	"context"

	"github.com/cloudwego/eino/components/model"
)

// newReportModel 创建报告生成模型
func newReportModel(ctx context.Context) (model.ToolCallingChatModel, error) {
	return models.OpenAIForDeepSeekV3Quick(ctx)
}
