package plugin

// Plugin 插件接口
type Plugin interface {
	Name() string
	Init() error
	Start() error
	Stop() error
}

// Info 插件信息
type Info struct {
	Name    string
	Version string
	Author  string
}

// Manager 插件管理器
type Manager struct {
	plugins map[string]Plugin
}

// NewManager 创建管理器
func NewManager() *Manager {
	return &Manager{
		plugins: make(map[string]Plugin),
	}
}

// Register 注册插件
func (m *Manager) Register(p Plugin) error {
	if err := p.Init(); err != nil {
		return err
	}
	m.plugins[p.Name()] = p
	return nil
}

// StartAll 启动所有插件
func (m *Manager) StartAll() error {
	for _, p := range m.plugins {
		if err := p.Start(); err != nil {
			return err
		}
	}
	return nil
}
