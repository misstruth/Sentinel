package event

import (
	"context"

	v1 "SuperBizAgent/api/event/v1"
	"SuperBizAgent/internal/model"
	"SuperBizAgent/internal/repository"
)

// UpdateStatus 更新事件状态
func (c *Controller) UpdateStatus(ctx context.Context, req *v1.UpdateStatusReq) (*v1.UpdateStatusRes, error) {
	repo := repository.NewEventRepository()
	err := repo.UpdateStatus(req.ID, model.EventStatus(req.Status))
	if err != nil {
		return nil, err
	}
	return &v1.UpdateStatusRes{}, nil
}
