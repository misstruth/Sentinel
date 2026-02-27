package report

import (
	"context"

	v1 "SuperBizAgent/api/report/v1"
)

// Delete 删除报告
func (c *Controller) Delete(ctx context.Context, req *v1.DeleteReq) (*v1.DeleteRes, error) {
	err := c.service.Delete(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &v1.DeleteRes{
		Success: true,
	}, nil
}
