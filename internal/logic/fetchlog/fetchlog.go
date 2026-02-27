package fetchlog

import (
	"context"

	"SuperBizAgent/internal/model"
	"SuperBizAgent/internal/repository"
)

// Service 抓取日志服务
type Service struct {
	repo *repository.FetchLogRepository
}

// NewService 创建服务
func NewService(repo *repository.FetchLogRepository) *Service {
	return &Service{repo: repo}
}

// List 获取日志列表
func (s *Service) List(ctx context.Context, subID uint, page, size int) ([]model.FetchLog, int64, error) {
	return s.repo.ListBySubscriptionPaged(subID, page, size)
}

// GetStats 获取统计信息
func (s *Service) GetStats(ctx context.Context, subID uint) (*Stats, error) {
	logs, err := s.repo.ListBySubscription(subID, 100)
	if err != nil {
		return nil, err
	}
	return s.calcStats(logs), nil
}

// Stats 统计信息
type Stats struct {
	TotalFetches   int `json:"total_fetches"`
	SuccessCount   int `json:"success_count"`
	FailedCount    int `json:"failed_count"`
	TotalEvents    int `json:"total_events"`
	AvgDurationMs  int `json:"avg_duration_ms"`
}

// calcStats 计算统计信息
func (s *Service) calcStats(logs []model.FetchLog) *Stats {
	stats := &Stats{TotalFetches: len(logs)}
	if len(logs) == 0 {
		return stats
	}

	var totalDuration int
	for _, l := range logs {
		if l.Status == model.FetchStatusSuccess {
			stats.SuccessCount++
		} else {
			stats.FailedCount++
		}
		stats.TotalEvents += l.EventCount
		totalDuration += l.Duration
	}
	stats.AvgDurationMs = totalDuration / len(logs)
	return stats
}
