package retriever

import "time"

// Config RAG 系统配置
type Config struct {
	// 混合检索配置
	HybridSearch HybridSearchConfig `json:"hybrid_search"`

	// Query 改写配置
	QueryRewrite QueryRewriteConfig `json:"query_rewrite"`

	// 重排序配置
	Rerank RerankConfig `json:"rerank"`

	// 缓存配置
	Cache CacheConfig `json:"cache"`

	// 降级配置
	Fallback FallbackConfig `json:"fallback"`
}

type HybridSearchConfig struct {
	Enabled      bool    `json:"enabled"`
	DenseWeight  float64 `json:"dense_weight"`  // 密集向量权重 0.7
	SparseWeight float64 `json:"sparse_weight"` // 稀疏向量权重 0.3
	TopK         int     `json:"top_k"`         // 召回数量 20
}

type QueryRewriteConfig struct {
	Enabled       bool `json:"enabled"`
	NumRewrites   int  `json:"num_rewrites"`   // 改写数量 3
	UseHyDE       bool `json:"use_hyde"`       // 是否使用 HyDE
}

type RerankConfig struct {
	Enabled         bool    `json:"enabled"`
	CoarseTopK      int     `json:"coarse_top_k"`      // 粗排 TopK 30
	FineTopK        int     `json:"fine_top_k"`        // 精排 TopK 10
	FinalTopK       int     `json:"final_top_k"`       // 最终 TopK 3
	ScoreThreshold  float64 `json:"score_threshold"`   // 分数阈值 0.6
	RerankURL       string  `json:"rerank_url"`        // Reranker 服务地址
	Timeout         int     `json:"timeout"`           // 超时时间(ms)
}

type CacheConfig struct {
	Enabled bool          `json:"enabled"`
	TTL     time.Duration `json:"ttl"` // 缓存过期时间 1h
	MaxSize int           `json:"max_size"`
}

type FallbackConfig struct {
	Enabled bool `json:"enabled"`
	Timeout int  `json:"timeout"` // 主检索超时时间(ms)
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		HybridSearch: HybridSearchConfig{
			Enabled:      true,
			DenseWeight:  0.7,
			SparseWeight: 0.3,
			TopK:         20,
		},
		QueryRewrite: QueryRewriteConfig{
			Enabled:     true,
			NumRewrites: 3,
			UseHyDE:     false,
		},
		Rerank: RerankConfig{
			Enabled:        true,
			CoarseTopK:     30,
			FineTopK:       10,
			FinalTopK:      3,
			ScoreThreshold: 0.6,
			Timeout:        3000,
		},
		Cache: CacheConfig{
			Enabled: true,
			TTL:     time.Hour,
			MaxSize: 10000,
		},
		Fallback: FallbackConfig{
			Enabled: true,
			Timeout: 5000,
		},
	}
}
