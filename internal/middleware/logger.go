package middleware

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Logger 日志中间件
func Logger(r *ghttp.Request) {
	start := time.Now()
	r.Middleware.Next()
	duration := time.Since(start)

	g.Log().Infof(r.Context(),
		"%s %s %d %v",
		r.Method,
		r.URL.Path,
		r.Response.Status,
		duration,
	)
}
