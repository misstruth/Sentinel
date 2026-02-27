package event

import (
	"context"

	v1 "SuperBizAgent/api/event/v1"
)

// DeleteAll 删除所有事件
func (c *Controller) DeleteAll(ctx context.Context, req *v1.DeleteAllReq) (*v1.DeleteAllRes, error) {
	err := c.service.DeleteAll(ctx)
	return &v1.DeleteAllRes{}, err
}
