package database

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/components/document"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"

	"SuperBizAgent/internal/ai/indexer"
	"SuperBizAgent/internal/ai/loader"
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

// SeedDocumentsToMilvus 索引文档到 Milvus 向量数据库
func SeedDocumentsToMilvus(ctx context.Context, docsPath string) error {
	idx, err := indexer.NewMilvusIndexer(ctx)
	if err != nil {
		return fmt.Errorf("创建索引器失败: %w", err)
	}

	fileLoader, err := loader.NewFileLoader(ctx)
	if err != nil {
		return fmt.Errorf("创建加载器失败: %w", err)
	}

	count := 0
	err = filepath.WalkDir(docsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
			return err
		}

		docs, err := fileLoader.Load(ctx, document.Source{URI: path})
		if err != nil {
			g.Log().Warningf(ctx, "加载文档失败 %s: %v", path, err)
			return nil
		}

		if _, err = idx.Store(ctx, docs); err != nil {
			g.Log().Warningf(ctx, "索引文档失败 %s: %v", path, err)
			return nil
		}

		count += len(docs)
		g.Log().Infof(ctx, "已索引文档: %s (%d 片段)", path, len(docs))
		return nil
	})

	if err != nil {
		return fmt.Errorf("遍历文档目录失败: %w", err)
	}

	g.Log().Infof(ctx, "seed: indexed %d document chunks to Milvus", count)
	return nil
}
