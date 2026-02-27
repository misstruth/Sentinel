package lock

import (
	"context"
	"sync"
)

// Lock 分布式锁接口
type Lock interface {
	Acquire(ctx context.Context) error
	Release(ctx context.Context) error
}

// MemoryLock 内存锁实现
type MemoryLock struct {
	key string
	mu  sync.Mutex
}

// NewMemoryLock 创建内存锁
func NewMemoryLock(key string) *MemoryLock {
	return &MemoryLock{key: key}
}

// Acquire 获取锁
func (l *MemoryLock) Acquire(ctx context.Context) error {
	l.mu.Lock()
	return nil
}

// Release 释放锁
func (l *MemoryLock) Release(ctx context.Context) error {
	l.mu.Unlock()
	return nil
}
