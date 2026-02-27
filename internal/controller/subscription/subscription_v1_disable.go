package subscription

import (
	"context"

	v1 "SuperBizAgent/api/subscription/v1"
	"SuperBizAgent/internal/logic/subscription"
)

// Disable 禁用订阅
func (c *ControllerV1) Disable(ctx context.Context, req *v1.DisableReq) (*v1.DisableRes, error) {
	svc := subscription.NewService()
	if err := svc.Disable(ctx, req.ID); err != nil {
		return nil, err
	}
	return &v1.DisableRes{}, nil
}
