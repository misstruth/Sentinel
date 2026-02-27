package hotreload

import (
	"sync"
)

// Watcher 配置监听器
type Watcher struct {
	path      string
	callbacks []func()
	mu        sync.RWMutex
}

// NewWatcher 创建监听器
func NewWatcher(path string) *Watcher {
	return &Watcher{
		path:      path,
		callbacks: make([]func(), 0),
	}
}

// OnChange 注册变更回调
func (w *Watcher) OnChange(fn func()) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.callbacks = append(w.callbacks, fn)
}

// Reload 触发重载
func (w *Watcher) Reload() {
	w.mu.RLock()
	defer w.mu.RUnlock()
	for _, fn := range w.callbacks {
		fn()
	}
}
