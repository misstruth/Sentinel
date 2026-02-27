package fetchlog

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

// Controller 抓取日志控制器
type Controller struct{}

// New 创建控制器
func New() *Controller {
	return &Controller{}
}

// List 获取日志列表
func (c *Controller) List(r *ghttp.Request) {
	subID := r.Get("subscription_id").Uint()
	page := r.Get("page", 1).Int()
	size := r.Get("size", 20).Int()

	r.Response.WriteJson(map[string]interface{}{
		"subscription_id": subID,
		"page":            page,
		"size":            size,
	})
}
