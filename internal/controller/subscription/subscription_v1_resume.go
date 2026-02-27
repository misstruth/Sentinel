package subscription

import (
	"context"

	v1 "SuperBizAgent/api/subscription/v1"
	"SuperBizAgent/internal/logic/subscription"
)

// Resume 恢复订阅
func (c *ControllerV1) Resume(ctx context.Context, req *v1.ResumeReq) (*v1.ResumeRes, error) {
	svc := subscription.NewService()
	if err := svc.Resume(ctx, req.ID); err != nil {
		return nil, err
	}
	return &v1.ResumeRes{}, nil
}
