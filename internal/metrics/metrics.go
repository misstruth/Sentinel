package metrics

import (
	"runtime"
	"sync"
	"time"
)

// Metrics 性能指标
type Metrics struct {
	RequestCount  int64
	ErrorCount    int64
	AvgLatency    float64
	MemoryUsage   uint64
	GoroutineNum  int
	LastUpdated   time.Time
}

// Collector 指标收集器
type Collector struct {
	metrics *Metrics
	mu      sync.RWMutex
}

// NewCollector 创建收集器
func NewCollector() *Collector {
	return &Collector{
		metrics: &Metrics{},
	}
}

// Collect 收集指标
func (c *Collector) Collect() *Metrics {
	c.mu.Lock()
	defer c.mu.Unlock()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	c.metrics.MemoryUsage = m.Alloc
	c.metrics.GoroutineNum = runtime.NumGoroutine()
	c.metrics.LastUpdated = time.Now()

	return c.metrics
}

// IncRequest 增加请求计数
func (c *Collector) IncRequest() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.metrics.RequestCount++
}
