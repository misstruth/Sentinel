package retry

import (
	"context"
	"time"
)

// Config 重试配置
type Config struct {
	MaxRetries int
	Delay      time.Duration
	MaxDelay   time.Duration
}

// Retryer 重试器
type Retryer struct {
	config *Config
}

// NewRetryer 创建重试器
func NewRetryer(cfg *Config) *Retryer {
	return &Retryer{config: cfg}
}

// Do 执行重试
func (r *Retryer) Do(ctx context.Context, fn func() error) error {
	var err error
	for i := 0; i <= r.config.MaxRetries; i++ {
		if err = fn(); err == nil {
			return nil
		}
		if i < r.config.MaxRetries {
			time.Sleep(r.config.Delay)
		}
	}
	return err
}
