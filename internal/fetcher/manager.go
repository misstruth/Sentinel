package fetcher

import (
	"context"
	"sync"

	"SuperBizAgent/internal/model"
)

// Manager Fetcher 管理器
type Manager struct {
	fetchers map[model.SourceType]Fetcher
	mu       sync.RWMutex
}

// NewManager 创建管理器
func NewManager() *Manager {
	return &Manager{
		fetchers: make(map[model.SourceType]Fetcher),
	}
}

// Register 注册 Fetcher
func (m *Manager) Register(f Fetcher) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.fetchers[f.Type()] = f
}

// Get 获取 Fetcher
func (m *Manager) Get(t model.SourceType) Fetcher {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.fetchers[t]
}

// Fetch 执行抓取
func (m *Manager) Fetch(ctx context.Context, sub *model.Subscription) ([]model.SecurityEvent, error) {
	f := m.Get(sub.SourceType)
	if f == nil {
		return nil, nil
	}
	return f.Fetch(ctx, sub)
}
