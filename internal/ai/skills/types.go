package skills

// Skill 定义
type Skill struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Category    string         `json:"category"`
	Enabled     bool           `json:"enabled"`
	Tools       []string       `json:"tools"`
	Prompt      string         `json:"-"`
	Params      []SkillParam   `json:"params"`
}

// SkillParam 参数定义
type SkillParam struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

// ExecuteRequest 执行请求
type ExecuteRequest struct {
	SkillID string         `json:"skill_id"`
	Params  map[string]any `json:"params"`
}

// ExecuteResult 执行结果
type ExecuteResult struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}
