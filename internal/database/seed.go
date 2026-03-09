package database

import (
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"

	"SuperBizAgent/internal/model"
)

// defaultSource 默认订阅源定义
type defaultSource struct {
	Name        string
	Type        model.SourceType
	URL         string
	Cron        string
	Description string
	Enabled     bool
}

// SeedDefaultSubscriptions 首次部署时写入默认订阅源
// 仅当 subscriptions 表为空时执行
func SeedDefaultSubscriptions() error {
	var count int64
	if err := DB.Model(&model.Subscription{}).Count(&count).Error; err != nil {
		return fmt.Errorf("check subscriptions count: %w", err)
	}
	if count > 0 {
		return nil
	}

	sources := getDefaultSources()
	if len(sources) == 0 {
		return nil
	}

	var subs []model.Subscription
	for _, src := range sources {
		status := model.StatusActive
		if !src.Enabled {
			status = model.StatusDisabled
		}
		subs = append(subs, model.Subscription{
			Name:        src.Name,
			Description: src.Description,
			SourceType:  src.Type,
			SourceURL:   src.URL,
			Status:      status,
			CronExpr:    src.Cron,
			AuthType:    "none",
		})
	}

	if err := DB.Create(&subs).Error; err != nil {
		return fmt.Errorf("seed subscriptions: %w", err)
	}

	ctx := gctx.New()
	g.Log().Infof(ctx, "seed: inserted %d default subscriptions", len(subs))
	return nil
}
