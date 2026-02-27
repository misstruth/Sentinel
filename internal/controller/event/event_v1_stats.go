package event

import (
	"context"

	v1 "SuperBizAgent/api/event/v1"
	"SuperBizAgent/internal/repository"
)

// Stats 获取事件统计
func (c *Controller) Stats(ctx context.Context, req *v1.StatsReq) (*v1.StatsRes, error) {
	repo := repository.NewEventRepository()
	total, bySeverity, byStatus, err := repo.GetStats()
	if err != nil {
		return nil, err
	}
	todayTotal, todayCritical, todayHigh := repo.GetTodayStats()

	return &v1.StatsRes{
		Total:         total,
		TodayCount:    todayTotal,
		CriticalCount: bySeverity["critical"] + todayCritical,
		HighCount:     bySeverity["high"] + todayHigh,
		BySeverity:    bySeverity,
		ByStatus:      byStatus,
	}, nil
}
