package event

import (
	"context"

	v1 "SuperBizAgent/api/event/v1"
	"SuperBizAgent/internal/model"
	"SuperBizAgent/internal/repository"
)

// BatchUpdateStatus 批量更新事件状态
func (c *Controller) BatchUpdateStatus(ctx context.Context, req *v1.BatchUpdateStatusReq) (*v1.BatchUpdateStatusRes, error) {
	repo := repository.NewEventRepository()
	err := repo.BatchUpdateStatus(req.IDs, model.EventStatus(req.Status))
	if err != nil {
		return nil, err
	}
	return &v1.BatchUpdateStatusRes{}, nil
}
