package subscription

import (
	"context"

	v1 "SuperBizAgent/api/subscription/v1"
	"SuperBizAgent/internal/repository"
)

// FetchLogs 获取抓取日志
func (c *ControllerV1) FetchLogs(ctx context.Context, req *v1.FetchLogsReq) (*v1.FetchLogsRes, error) {
	repo := repository.NewFetchLogRepository()

	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	logs, total, err := repo.ListBySubscriptionPaged(req.ID, page, pageSize)
	if err != nil {
		return nil, err
	}

	list := make([]*v1.FetchLogItem, len(logs))
	for i, log := range logs {
		list[i] = &v1.FetchLogItem{
			ID:             log.ID,
			SubscriptionID: log.SubscriptionID,
			Status:         string(log.Status),
			EventCount:     log.EventCount,
			ErrorMsg:       log.ErrorMsg,
			Duration:       log.Duration,
			CreatedAt:      log.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return &v1.FetchLogsRes{
		List:  list,
		Total: total,
		Page:  page,
		Size:  pageSize,
	}, nil
}
