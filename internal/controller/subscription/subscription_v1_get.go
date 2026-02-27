package subscription

import (
	"context"

	v1 "SuperBizAgent/api/subscription/v1"
	"SuperBizAgent/internal/logic/subscription"
)

// Get 获取订阅详情
func (c *ControllerV1) Get(ctx context.Context, req *v1.GetReq) (*v1.GetRes, error) {
	svc := subscription.NewService()

	sub, err := svc.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	lastFetchAt := ""
	if sub.LastFetchAt != nil {
		lastFetchAt = sub.LastFetchAt.Format("2006-01-02 15:04:05")
	}

	return &v1.GetRes{
		SubscriptionItem: &v1.SubscriptionItem{
			ID:          sub.ID,
			Name:        sub.Name,
			Description: sub.Description,
			SourceType:  string(sub.SourceType),
			SourceURL:   sub.SourceURL,
			Status:      string(sub.Status),
			CronExpr:    sub.CronExpr,
			TotalEvents: sub.TotalEvents,
			LastFetchAt: lastFetchAt,
			CreatedAt:   sub.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}
