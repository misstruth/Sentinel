package subscription

import (
	"context"

	v1 "SuperBizAgent/api/subscription/v1"
	"SuperBizAgent/internal/logic/subscription"
	"SuperBizAgent/internal/model"
)

// List 订阅列表
func (c *ControllerV1) List(ctx context.Context, req *v1.ListReq) (*v1.ListRes, error) {
	svc := subscription.NewService()

	query := &subscription.ListQuery{
		Page:       req.Page,
		PageSize:   req.PageSize,
		SourceType: model.SourceType(req.SourceType),
		Status:     model.SubscriptionStatus(req.Status),
		Keyword:    req.Keyword,
	}

	result, err := svc.List(ctx, query)
	if err != nil {
		return nil, err
	}

	items := make([]*v1.SubscriptionItem, 0, len(result.Items))
	for _, sub := range result.Items {
		lastFetchAt := ""
		if sub.LastFetchAt != nil {
			lastFetchAt = sub.LastFetchAt.Format("2006-01-02 15:04:05")
		}
		items = append(items, &v1.SubscriptionItem{
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
		})
	}

	return &v1.ListRes{
		Items: items,
		Total: result.Total,
		Page:  result.Page,
		Size:  result.Size,
	}, nil
}
