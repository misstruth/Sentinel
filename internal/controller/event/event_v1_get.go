package event

import (
	"context"

	v1 "SuperBizAgent/api/event/v1"
)

// Get 获取事件详情
func (c *Controller) Get(ctx context.Context, req *v1.GetReq) (*v1.GetRes, error) {
	e, err := c.service.Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	// 模拟受影响资产数（基于风险评分）
	affectedAssets := 0
	if e.RiskScore >= 80 {
		affectedAssets = 5 + int(e.ID%3)
	} else if e.RiskScore >= 60 {
		affectedAssets = 2 + int(e.ID%2)
	} else if e.RiskScore > 0 {
		affectedAssets = 1
	}

	return &v1.GetRes{
		ID:             e.ID,
		Title:          e.Title,
		Description:    e.Description,
		Severity:       string(e.Severity),
		Status:         string(e.Status),
		CVEID:          e.CVEID,
		CVSSScore:      e.CVSSScore,
		SourceURL:      e.SourceURL,
		EventTime:      e.EventTime.Format("2006-01-02 15:04:05"),
		CreatedAt:      e.CreatedAt.Format("2006-01-02 15:04:05"),
		RiskScore:      e.RiskScore,
		Recommendation: e.Recommendation,
		AffectedAssets: affectedAssets,
	}, nil
}
