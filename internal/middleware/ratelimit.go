package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
)

var limiter = newRateLimiter(100, time.Minute)

type rateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	rl := &rateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
	go rl.cleanup()
	return rl
}

// RateLimit 限流中间件
func RateLimit(r *ghttp.Request) {
	ip := r.GetClientIp()
	if !limiter.allow(ip) {
		r.Response.WriteStatus(http.StatusTooManyRequests)
		r.Response.WriteJson(map[string]string{"error": "rate limit exceeded"})
		return
	}
	r.Middleware.Next()
}

func (l *rateLimiter) allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	cutoff := now.Add(-l.window)

	// 过滤过期请求
	valid := l.requests[key][:0]
	for _, t := range l.requests[key] {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= l.limit {
		l.requests[key] = valid
		return false
	}
	l.requests[key] = append(valid, now)
	return true
}

func (l *rateLimiter) cleanup() {
	for {
		time.Sleep(5 * time.Minute)
		l.mu.Lock()
		cutoff := time.Now().Add(-l.window)
		for k, times := range l.requests {
			valid := times[:0]
			for _, t := range times {
				if t.After(cutoff) {
					valid = append(valid, t)
				}
			}
			if len(valid) == 0 {
				delete(l.requests, k)
			} else {
				l.requests[k] = valid
			}
		}
		l.mu.Unlock()
	}
}
