package supervisor

// AgentType 子 Agent 类型
type AgentType string

const (
	AgentChat       AgentType = "chat"
	AgentEvent      AgentType = "event"
	AgentReport     AgentType = "report"
	AgentRisk       AgentType = "risk"
	AgentPlan       AgentType = "plan"
)

// Task 任务
type Task struct {
	ID      string
	Query   string
	Type    AgentType
	Params  map[string]any
}

// Result 结果
type Result struct {
	TaskID  string
	Agent   AgentType
	Content string
	Error   error
}

// StreamCallback 流式回调
type StreamCallback func(agent AgentType, chunk string)

// SupervisorInput Graph 输入
type SupervisorInput struct {
	Query    string
	Callback StreamCallback
}

// SupervisorOutput Graph 输出
type SupervisorOutput struct {
	TaskID  string
	Agent   AgentType
	Content string
	Error   error
}
