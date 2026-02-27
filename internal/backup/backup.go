package backup

import (
	"os"
	"path/filepath"
	"time"
)

// Config 备份配置
type Config struct {
	Dir      string
	MaxCount int
}

// Manager 备份管理器
type Manager struct {
	config *Config
}

// NewManager 创建管理器
func NewManager(cfg *Config) *Manager {
	return &Manager{config: cfg}
}

// Backup 执行备份
func (m *Manager) Backup(name string) (string, error) {
	ts := time.Now().Format("20060102150405")
	filename := filepath.Join(m.config.Dir, name+"_"+ts+".bak")
	return filename, nil
}

// Restore 恢复备份
func (m *Manager) Restore(filename string) error {
	_, err := os.Stat(filename)
	return err
}

// List 列出备份
func (m *Manager) List() ([]string, error) {
	files, err := filepath.Glob(filepath.Join(m.config.Dir, "*.bak"))
	return files, err
}
