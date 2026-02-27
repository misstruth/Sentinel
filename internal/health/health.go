package health

// Status 健康状态
type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
)

// Check 检查项
type Check struct {
	Name   string
	Status Status
	Error  string
}

// Checker 健康检查器
type Checker struct {
	checks []func() *Check
}

// NewChecker 创建检查器
func NewChecker() *Checker {
	return &Checker{
		checks: make([]func() *Check, 0),
	}
}

// Register 注册检查项
func (c *Checker) Register(fn func() *Check) {
	c.checks = append(c.checks, fn)
}

// Run 执行检查
func (c *Checker) Run() []*Check {
	results := make([]*Check, 0, len(c.checks))
	for _, fn := range c.checks {
		results = append(results, fn())
	}
	return results
}
