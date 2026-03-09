package retriever

import (
	"context"
	"time"
)

// Document 检索文档结构
type Document struct {
	ID       string                 `json:"id"`
	Content  string                 `json:"content"`
	Score    float64                `json:"score"`
	Metadata map[string]interface{} `json:"metadata"`
}

// QueryAnalysis 查询分析结果
type QueryAnalysis struct {
	OriginalQuery string   `json:"original_query"`
	Keywords      []string `json:"keywords"`
	Complexity    string   `json:"complexity"` // "simple" | "medium" | "complex"
	TopK          int      `json:"top_k"`
	RewriteQueries []string `json:"rewrite_queries"`
}

// RetrievalResult 检索结果
type RetrievalResult struct {
	Documents []Document `json:"documents"`
	Latency   time.Duration `json:"latency"`
	Source    string     `json:"source"` // "vector" | "hybrid" | "cache"
}

// Retriever 检索器接口
type Retriever interface {
	Retrieve(ctx context.Context, query string) ([]Document, error)
}

// QueryAnalyzer 查询分析器接口
type QueryAnalyzer interface {
	Analyze(ctx context.Context, query string) (*QueryAnalysis, error)
}

// Reranker 重排序器接口
type Reranker interface {
	Rerank(ctx context.Context, query string, docs []Document) ([]Document, error)
}

// Cache 缓存接口
type Cache interface {
	Get(ctx context.Context, key string) ([]Document, bool, error)
	Set(ctx context.Context, key string, docs []Document, ttl time.Duration) error
}
