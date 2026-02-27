package subscription

import (
	"context"

	v1 "SuperBizAgent/api/subscription/v1"
	"SuperBizAgent/internal/logic/subscription"
	"SuperBizAgent/internal/service"
)

// Fetch 手动触发抓取
func (c *ControllerV1) Fetch(ctx context.Context, req *v1.FetchReq) (*v1.FetchRes, error) {
	svc := subscription.NewService()
	sub, err := svc.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	result, err := service.FetchSubscription(sub)
	if err != nil {
		return &v1.FetchRes{
			Duration: result.Duration,
			Message:  "抓取失败: " + err.Error(),
		}, nil
	}

	return &v1.FetchRes{
		FetchedCount: result.FetchedCount,
		NewCount:     result.NewCount,
		TotalEvents:  result.TotalEvents,
		Duration:     result.Duration,
		Message:      "抓取成功",
	}, nil
}
