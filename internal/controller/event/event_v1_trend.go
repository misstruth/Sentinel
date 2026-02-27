package event

import (
	"context"

	v1 "SuperBizAgent/api/event/v1"
	"SuperBizAgent/internal/repository"
)

// Trend 获取事件趋势
func (c *Controller) Trend(ctx context.Context, req *v1.TrendReq) (*v1.TrendRes, error) {
	repo := repository.NewEventRepository()

	days := req.Days
	if days < 1 {
		days = 7
	}
	if days > 30 {
		days = 30
	}

	data, err := repo.GetTrend(days)
	if err != nil {
		return nil, err
	}

	list := make([]v1.TrendItem, len(data))
	for i, d := range data {
		list[i] = v1.TrendItem{
			Date:     d["date"].(string),
			Total:    d["total"].(int64),
			Critical: d["critical"].(int64),
			High:     d["high"].(int64),
			Medium:   d["medium"].(int64),
			Low:      d["low"].(int64),
			Info:     d["info"].(int64),
		}
	}

	return &v1.TrendRes{List: list}, nil
}
