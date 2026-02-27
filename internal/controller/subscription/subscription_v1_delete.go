package subscription

import (
	"context"

	v1 "SuperBizAgent/api/subscription/v1"
	"SuperBizAgent/internal/logic/subscription"
)

// Delete 删除订阅
func (c *ControllerV1) Delete(ctx context.Context, req *v1.DeleteReq) (*v1.DeleteRes, error) {
	svc := subscription.NewService()

	if err := svc.Delete(ctx, req.ID); err != nil {
		return nil, err
	}

	return &v1.DeleteRes{}, nil
}
