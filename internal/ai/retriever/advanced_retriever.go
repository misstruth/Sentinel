package retriever

import (
	"context"
	"log"
	"time"

	"github.com/cloudwego/eino/components/model"
)

// AdvancedRetriever 高级检索器（整合所有功能）
type AdvancedRetriever struct {
	config        *Config
	baseRetriever *HybridRetriever
	queryRewriter *QueryRewriter
	reranker      *MultiStageReranker
	cache         Cache
	fusion        *RRFFusion
	analyzer      *QueryAnalyzerImpl
}

func NewAdvancedRetriever(ctx context.Context, llm model.ChatModel, config *Config) (*AdvancedRetriever, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// 初始化基础检索器
	baseRetriever, err := NewHybridRetriever(ctx, config.HybridSearch.TopK)
	if err != nil {
		return nil, err
	}

	// 初始化 Query 改写器
	var queryRewriter *QueryRewriter
	if config.QueryRewrite.Enabled {
		queryRewriter = NewQueryRewriter(llm, config.QueryRewrite)
	}

	// 初始化重排序器
	var reranker *MultiStageReranker
	if config.Rerank.Enabled {
		reranker = NewMultiStageReranker(config.Rerank)
	}

	// 初始化缓存（可选）
	var cache Cache
	if config.Cache.Enabled {
		// 这里需要 Redis 地址，暂时留空，实际使用时从配置读取
		// cache, _ = NewRedisCache("localhost:6379", config.Cache.TTL)
	}

	return &AdvancedRetriever{
		config:        config,
		baseRetriever: baseRetriever,
		queryRewriter: queryRewriter,
		reranker:      reranker,
		cache:         cache,
		fusion:        NewRRFFusion(),
		analyzer:      NewQueryAnalyzer(llm),
	}, nil
}

func (a *AdvancedRetriever) Retrieve(ctx context.Context, query string) ([]Document, error) {
	startTime := time.Now()

	// 1. 检查缓存
	if a.cache != nil {
		if docs, hit, err := a.cache.Get(ctx, query); err == nil && hit {
			log.Printf("[Cache Hit] query=%s, latency=%v", query, time.Since(startTime))
			return docs, nil
		}
	}

	// 2. Query 改写
	queries := []string{query}
	if a.queryRewriter != nil {
		rewrites, err := a.queryRewriter.Rewrite(ctx, query)
		if err == nil {
			queries = rewrites
		}
	}

	// 3. 多路召回
	allResults := [][]Document{}
	for _, q := range queries {
		docs, err := a.baseRetriever.Retrieve(ctx, q)
		if err != nil {
			log.Printf("[Retrieve Error] query=%s, err=%v", q, err)
			continue
		}
		allResults = append(allResults, docs)
	}

	if len(allResults) == 0 {
		return []Document{}, nil
	}

	// 4. 融合结果
	merged := a.fusion.Fuse(allResults)
	merged = Deduplicate(merged)

	// 5. 重排序
	if a.reranker != nil {
		reranked, err := a.reranker.Rerank(ctx, query, merged)
		if err == nil {
			merged = reranked
		}
	}

	// 6. 写入缓存
	if a.cache != nil {
		_ = a.cache.Set(ctx, query, merged, a.config.Cache.TTL)
	}

	log.Printf("[Retrieve Success] query=%s, results=%d, latency=%v", query, len(merged), time.Since(startTime))
	return merged, nil
}
