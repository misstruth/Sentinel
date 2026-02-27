package middleware

import (
	"net/http"
	"strings"

	"SuperBizAgent/internal/auth/jwt"

	"github.com/gogf/gf/v2/net/ghttp"
)

// Auth 认证中间件
func Auth(r *ghttp.Request) {
	token := r.GetHeader("Authorization")
	if token == "" {
		r.Response.WriteStatus(http.StatusUnauthorized)
		r.Response.WriteJson(map[string]string{
			"error": "missing authorization header",
		})
		return
	}

	token = strings.TrimPrefix(token, "Bearer ")
	if !validateToken(token) {
		r.Response.WriteStatus(http.StatusUnauthorized)
		r.Response.WriteJson(map[string]string{
			"error": "invalid token",
		})
		return
	}

	r.Middleware.Next()
}

// validateToken 验证 token
func validateToken(token string) bool {
	_, err := jwt.Parse(token)
	return err == nil
}
