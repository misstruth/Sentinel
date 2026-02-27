package service

import (
	"context"
	"log"
	"time"

	"SuperBizAgent/internal/fetcher"
	"SuperBizAgent/internal/model"
	"SuperBizAgent/internal/repository"
	"SuperBizAgent/internal/scheduler"
)

// FetchService 数据抓取服务
type FetchService struct {
	fetcherMgr *fetcher.Manager
	scheduler  *scheduler.Scheduler
	eventRepo  *repository.EventRepository
	logRepo    *repository.FetchLogRepository
}

// NewFetchService 创建抓取服务
func NewFetchService(eventRepo *repository.EventRepository, logRepo *repository.FetchLogRepository) *FetchService {
	svc := &FetchService{
		fetcherMgr: fetcher.NewManager(),
		eventRepo:  eventRepo,
		logRepo:    logRepo,
	}
	svc.scheduler = scheduler.New(svc.handleFetch)
	return svc
}

// RegisterFetchers 注册抓取器
func (s *FetchService) RegisterFetchers(githubToken, nvdKey string) {
	s.fetcherMgr.Register(fetcher.NewGitHubFetcher())
	s.fetcherMgr.Register(fetcher.NewRSSFetcher())
	s.fetcherMgr.Register(fetcher.NewNVDFetcher(nvdKey))
}

// Start 启动服务
func (s *FetchService) Start() {
	s.scheduler.Start()
}

// Stop 停止服务
func (s *FetchService) Stop() {
	s.scheduler.Stop()
}

// AddSubscription 添加订阅
func (s *FetchService) AddSubscription(sub *model.Subscription) {
	s.scheduler.Add(sub)
}

// RemoveSubscription 移除订阅
func (s *FetchService) RemoveSubscription(id uint) {
	s.scheduler.Remove(id)
}

// handleFetch 处理抓取任务
func (s *FetchService) handleFetch(ctx context.Context, sub *model.Subscription) error {
	start := time.Now()

	events, err := s.fetcherMgr.Fetch(ctx, sub)
	duration := time.Since(start)

	// 记录抓取日志
	fetchLog := &model.FetchLog{
		SubscriptionID: sub.ID,
		Duration:       int(duration.Milliseconds()),
	}

	if err != nil {
		fetchLog.Status = model.FetchStatusFailed
		fetchLog.ErrorMsg = err.Error()
		s.logRepo.Create(fetchLog)
		return err
	}

	fetchLog.Status = model.FetchStatusSuccess
	fetchLog.EventCount = len(events)
	s.logRepo.Create(fetchLog)

	// 保存事件
	s.saveEvents(events)
	return nil
}

// saveEvents 保存事件（去重）
func (s *FetchService) saveEvents(events []model.SecurityEvent) {
	for _, e := range events {
		if !s.eventRepo.ExistsByHash(e.UniqueHash) {
			if err := s.eventRepo.Create(&e); err != nil {
				log.Printf("save event error: %v", err)
			}
		}
	}
}
