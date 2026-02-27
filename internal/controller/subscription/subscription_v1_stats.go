package subscription

import (
	"context"

	v1 "SuperBizAgent/api/subscription/v1"
	"SuperBizAgent/internal/logic/subscription"
	"SuperBizAgent/internal/repository"
)

// Stats 获取订阅统计
func (c *ControllerV1) Stats(ctx context.Context, req *v1.StatsReq) (*v1.StatsRes, error) {
	svc := subscription.NewService()
	sub, err := svc.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	logRepo := repository.NewFetchLogRepository()
	total, success, failed, avgDuration, _ := logRepo.GetStats(req.ID)

	return &v1.StatsRes{
		TotalFetches:  total,
		SuccessCount:  success,
		FailedCount:   failed,
		TotalEvents:   sub.TotalEvents,
		AvgDurationMs: avgDuration,
	}, nil
}
