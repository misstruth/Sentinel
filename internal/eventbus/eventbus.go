package eventbus

import (
	"sync"
)

// Event 事件
type Event struct {
	Name string
	Data interface{}
}

// Handler 事件处理器
type Handler func(e *Event)

// Bus 事件总线
type Bus struct {
	handlers map[string][]Handler
	mu       sync.RWMutex
}

// NewBus 创建事件总线
func NewBus() *Bus {
	return &Bus{
		handlers: make(map[string][]Handler),
	}
}

// Subscribe 订阅事件
func (b *Bus) Subscribe(name string, h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[name] = append(b.handlers[name], h)
}

// Publish 发布事件
func (b *Bus) Publish(e *Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, h := range b.handlers[e.Name] {
		go h(e)
	}
}
