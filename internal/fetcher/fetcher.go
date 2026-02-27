package fetcher

import (
	"context"
	"time"

	"SuperBizAgent/internal/model"
)

// Fetcher 数据抓取器接口
type Fetcher interface {
	Fetch(ctx context.Context, sub *model.Subscription) ([]model.SecurityEvent, error)
	Type() model.SourceType
}

// Result 抓取结果
type Result struct {
	Events   []model.SecurityEvent
	Duration time.Duration
	Error    error
}
