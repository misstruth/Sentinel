package version

// Version API 版本
type Version string

const (
	V1 Version = "v1"
	V2 Version = "v2"
)

// Router 版本路由
type Router struct {
	handlers map[Version]interface{}
}

// NewRouter 创建路由
func NewRouter() *Router {
	return &Router{
		handlers: make(map[Version]interface{}),
	}
}

// Register 注册版本处理器
func (r *Router) Register(v Version, h interface{}) {
	r.handlers[v] = h
}

// Get 获取处理器
func (r *Router) Get(v Version) interface{} {
	return r.handlers[v]
}
