package retriever

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	retrievalLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "rag_retrieval_latency_ms",
			Help:    "RAG retrieval latency in milliseconds",
			Buckets: []float64{10, 50, 100, 200, 500, 1000, 2000},
		},
		[]string{"stage"},
	)

	cacheHitCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rag_cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"status"},
	)

	retrievalCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rag_retrieval_total",
			Help: "Total number of retrievals",
		},
		[]string{"status"},
	)
)

// RecordLatency 记录延迟
func RecordLatency(stage string, duration time.Duration) {
	retrievalLatency.WithLabelValues(stage).Observe(float64(duration.Milliseconds()))
}

// RecordCacheHit 记录缓存命中
func RecordCacheHit(hit bool) {
	status := "miss"
	if hit {
		status = "hit"
	}
	cacheHitCounter.WithLabelValues(status).Inc()
}

// RecordRetrieval 记录检索
func RecordRetrieval(success bool) {
	status := "failure"
	if success {
		status = "success"
	}
	retrievalCounter.WithLabelValues(status).Inc()
}
