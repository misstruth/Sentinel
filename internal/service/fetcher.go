package service

import (
	"context"
	"time"

	"SuperBizAgent/internal/database"
	"SuperBizAgent/internal/fetcher"
	"SuperBizAgent/internal/model"
)

var fetchManager *fetcher.Manager

func InitFetcher() {
	fetchManager = fetcher.NewManager()
	fetchManager.Register(fetcher.NewRSSFetcher())
	fetchManager.Register(fetcher.NewGitHubFetcher())
	go runFetchLoop()
}

func runFetchLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		fetchAllActive()
	}
}

func fetchAllActive() {
	db := database.GetDB()
	var subs []model.Subscription
	db.Where("status = ?", model.StatusActive).Find(&subs)

	// Worker pool限制并发数
	sem := make(chan struct{}, 10)
	for _, sub := range subs {
		sem <- struct{}{}
		go func(s model.Subscription) {
			defer func() { <-sem }()
			FetchSubscription(&s)
		}(sub)
	}
}

// FetchResult 抓取结果
type FetchResult struct {
	FetchedCount int   // 抓取到的数量
	NewCount     int   // 新增数量
	TotalEvents  int64 // 总事件数
	Duration     int   // 耗时ms
}

func FetchSubscription(sub *model.Subscription) (*FetchResult, error) {
	if fetchManager == nil {
		return &FetchResult{}, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	start := time.Now()
	events, err := fetchManager.Fetch(ctx, sub)
	duration := int(time.Since(start).Milliseconds())

	db := database.GetDB()
	log := &model.FetchLog{
		SubscriptionID: sub.ID,
		Duration:       duration,
		CreatedAt:      time.Now(),
	}

	if err != nil {
		log.Status = model.FetchStatusFailed
		log.ErrorMsg = err.Error()
		db.Create(log)
		return &FetchResult{Duration: duration}, err
	}

	log.Status = model.FetchStatusSuccess
	log.EventCount = len(events)
	db.Create(log)

	// 保存事件(去重)
	newCount := 0
	for _, e := range events {
		var count int64
		db.Model(&model.SecurityEvent{}).Where("unique_hash = ?", e.UniqueHash).Count(&count)
		if count == 0 {
			db.Create(&e)
			newCount++
		}
	}

	// 更新订阅统计 - 累计总事件数
	var totalEvents int64
	db.Model(&model.SecurityEvent{}).Where("subscription_id = ?", sub.ID).Count(&totalEvents)
	now := time.Now()
	db.Model(sub).Updates(map[string]interface{}{
		"last_fetch_at": now,
		"total_events":  totalEvents,
	})

	return &FetchResult{
		FetchedCount: len(events),
		NewCount:     newCount,
		TotalEvents:  totalEvents,
		Duration:     duration,
	}, nil
}
