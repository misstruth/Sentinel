package concurrency

// Limiter 并发限制器
type Limiter struct {
	sem chan struct{}
}

// NewLimiter 创建限制器
func NewLimiter(max int) *Limiter {
	return &Limiter{
		sem: make(chan struct{}, max),
	}
}

// Acquire 获取许可
func (l *Limiter) Acquire() {
	l.sem <- struct{}{}
}

// Release 释放许可
func (l *Limiter) Release() {
	<-l.sem
}
