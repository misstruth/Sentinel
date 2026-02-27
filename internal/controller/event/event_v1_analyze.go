package event

import (
	"context"

	v1 "SuperBizAgent/api/event/v1"
	"SuperBizAgent/internal/database"
	"SuperBizAgent/internal/model"
	"SuperBizAgent/internal/service"
)

// Analyze AI分析事件
func (c *Controller) Analyze(ctx context.Context, req *v1.AnalyzeReq) (*v1.AnalyzeRes, error) {
	if err := service.AnalyzeEvent(ctx, req.ID); err != nil {
		return nil, err
	}

	db := database.GetDB()
	var event model.SecurityEvent
	db.First(&event, req.ID)

	return &v1.AnalyzeRes{
		RiskScore:      event.RiskScore,
		Severity:       string(event.Severity),
		Recommendation: event.Recommendation,
	}, nil
}
