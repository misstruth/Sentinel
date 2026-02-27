package migrations

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Migration 迁移记录
type Migration struct {
	ID        uint      `gorm:"primaryKey"`
	Version   string    `gorm:"size:50;uniqueIndex"`
	Name      string    `gorm:"size:200"`
	AppliedAt time.Time `gorm:"autoCreateTime"`
}

// TableName 表名
func (Migration) TableName() string {
	return "migrations"
}

// Migrator 迁移管理器
type Migrator struct {
	db *gorm.DB
}

// NewMigrator 创建迁移管理器
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{db: db}
}

// Init 初始化迁移表
func (m *Migrator) Init() error {
	return m.db.AutoMigrate(&Migration{})
}

// IsApplied 检查迁移是否已应用
func (m *Migrator) IsApplied(version string) bool {
	var count int64
	m.db.Model(&Migration{}).Where("version = ?", version).Count(&count)
	return count > 0
}

// Apply 应用迁移
func (m *Migrator) Apply(version, name string, fn func(*gorm.DB) error) error {
	if m.IsApplied(version) {
		return nil
	}

	if err := fn(m.db); err != nil {
		return fmt.Errorf("迁移 %s 失败: %w", version, err)
	}

	return m.db.Create(&Migration{Version: version, Name: name}).Error
}
