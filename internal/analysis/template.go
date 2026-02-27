package analysis

// Template 分析模板
type Template struct {
	ID          uint
	Name        string
	Description string
	Prompt      string
	Category    string
	IsBuiltin   bool
}

// Manager 模板管理器
type Manager struct {
	templates map[uint]*Template
}

// NewManager 创建模板管理器
func NewManager() *Manager {
	return &Manager{
		templates: make(map[uint]*Template),
	}
}

// Add 添加模板
func (m *Manager) Add(tpl *Template) {
	m.templates[tpl.ID] = tpl
}

// Get 获取模板
func (m *Manager) Get(id uint) *Template {
	return m.templates[id]
}

// List 列出所有模板
func (m *Manager) List() []*Template {
	list := make([]*Template, 0, len(m.templates))
	for _, tpl := range m.templates {
		list = append(list, tpl)
	}
	return list
}

// Delete 删除模板
func (m *Manager) Delete(id uint) {
	delete(m.templates, id)
}
