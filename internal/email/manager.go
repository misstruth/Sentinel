package email

import (
	"sync"
)

// MailboxConfig 邮箱配置
type MailboxConfig struct {
	Name     string
	Host     string
	Port     int
	Username string
	Password string
	From     string
	UseTLS   bool
	IsActive bool
}

// Manager 多邮箱管理器
type Manager struct {
	mailboxes map[string]*MailboxConfig
	mu        sync.RWMutex
}

// NewManager 创建邮箱管理器
func NewManager() *Manager {
	return &Manager{
		mailboxes: make(map[string]*MailboxConfig),
	}
}

// Add 添加邮箱配置
func (m *Manager) Add(cfg *MailboxConfig) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.mailboxes[cfg.Name] = cfg
}

// Remove 移除邮箱配置
func (m *Manager) Remove(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.mailboxes, name)
}

// Get 获取邮箱配置
func (m *Manager) Get(name string) *MailboxConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.mailboxes[name]
}

// List 列出所有邮箱
func (m *Manager) List() []*MailboxConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	list := make([]*MailboxConfig, 0, len(m.mailboxes))
	for _, cfg := range m.mailboxes {
		list = append(list, cfg)
	}
	return list
}

// GetActive 获取活跃邮箱列表
func (m *Manager) GetActive() []*MailboxConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	list := make([]*MailboxConfig, 0)
	for _, cfg := range m.mailboxes {
		if cfg.IsActive {
			list = append(list, cfg)
		}
	}
	return list
}

// SetActive 设置邮箱激活状态
func (m *Manager) SetActive(name string, active bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if cfg, ok := m.mailboxes[name]; ok {
		cfg.IsActive = active
	}
}
