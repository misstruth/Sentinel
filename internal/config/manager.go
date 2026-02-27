package config

import (
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

// Manager 配置管理器
type Manager struct {
	mu       sync.RWMutex
	watchers []func()
}

var manager = &Manager{}

// GetManager 获取配置管理器
func GetManager() *Manager {
	return manager
}

// Get 获取配置值
func (m *Manager) Get(key string) interface{} {
	ctx := gctx.New()
	val, _ := g.Cfg().Get(ctx, key)
	return val.Interface()
}

// GetString 获取字符串配置
func (m *Manager) GetString(key string) string {
	ctx := gctx.New()
	val, _ := g.Cfg().Get(ctx, key)
	return val.String()
}

// Watch 注册配置变更监听器
func (m *Manager) Watch(fn func()) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.watchers = append(m.watchers, fn)
}

// NotifyChange 通知配置变更
func (m *Manager) NotifyChange() {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, fn := range m.watchers {
		go fn()
	}
}
