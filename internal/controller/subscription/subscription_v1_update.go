package subscription

import (
	"context"

	v1 "SuperBizAgent/api/subscription/v1"
	"SuperBizAgent/internal/logic/subscription"
	"SuperBizAgent/internal/model"
)

// Update 更新订阅
func (c *ControllerV1) Update(ctx context.Context, req *v1.UpdateReq) (*v1.UpdateRes, error) {
	svc := subscription.NewService()

	sub, err := svc.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		sub.Name = req.Name
	}
	if req.Description != "" {
		sub.Description = req.Description
	}
	if req.SourceURL != "" {
		sub.SourceURL = req.SourceURL
	}
	if req.Status != "" {
		sub.Status = model.SubscriptionStatus(req.Status)
	}
	if req.CronExpr != "" {
		sub.CronExpr = req.CronExpr
	}
	if req.Config != "" {
		sub.Config = req.Config
	}

	if err := svc.Update(ctx, sub); err != nil {
		return nil, err
	}

	return &v1.UpdateRes{}, nil
}
