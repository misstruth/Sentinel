package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// ReportType 报告类型
type ReportType string

const (
	ReportTypeDaily      ReportType = "daily"      // 日报
	ReportTypeWeekly     ReportType = "weekly"     // 周报
	ReportTypeMonthly    ReportType = "monthly"    // 月报
	ReportTypeCustom     ReportType = "custom"     // 自定义
	ReportTypeVulnAlert  ReportType = "vuln_alert" // 漏洞告警报告
	ReportTypeThreatBrief ReportType = "threat_brief" // 威胁简报
)

// ReportStatus 报告状态
type ReportStatus string

const (
	ReportStatusPending   ReportStatus = "pending"   // 待生成
	ReportStatusGenerating ReportStatus = "generating" // 生成中
	ReportStatusCompleted ReportStatus = "completed" // 已完成
	ReportStatusFailed    ReportStatus = "failed"    // 生成失败
)

// Report 安全报告模型
type Report struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"size:255;not null" json:"title"`
	Type      ReportType     `gorm:"size:20;index" json:"type"`
	Status    ReportStatus   `gorm:"size:20;default:pending" json:"status"`
	StartTime time.Time      `json:"start_time"`
	EndTime   time.Time      `json:"end_time"`
	Content   string         `gorm:"type:longtext" json:"content"`
	Summary   string         `gorm:"type:text" json:"summary"`

	// 模板配置
	TemplateID   uint   `gorm:"index" json:"template_id"`
	TemplateData string `gorm:"type:text" json:"template_data"` // JSON格式模板变量

	// 关联数据
	EventIDs       string `gorm:"type:text" json:"event_ids"`       // 关联事件ID列表(JSON)
	SubscriptionIDs string `gorm:"type:text" json:"subscription_ids"` // 关联订阅源ID列表(JSON)

	// 统计信息
	EventCount    int `gorm:"default:0" json:"event_count"`    // 事件总数
	CriticalCount int `gorm:"default:0" json:"critical_count"` // 严重事件数
	HighCount     int `gorm:"default:0" json:"high_count"`     // 高危事件数

	// 生成配置
	GeneratedBy string `gorm:"size:50" json:"generated_by"` // 生成方式: manual/scheduled/api
	ErrorMsg    string `gorm:"type:text" json:"error_msg"`  // 错误信息

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// TableName 指定表名
func (Report) TableName() string {
	return "reports"
}

// GetEventIDs 获取关联事件ID列表
func (r *Report) GetEventIDs() []uint {
	if r.EventIDs == "" {
		return nil
	}
	var ids []uint
	json.Unmarshal([]byte(r.EventIDs), &ids)
	return ids
}

// SetEventIDs 设置关联事件ID列表
func (r *Report) SetEventIDs(ids []uint) error {
	data, err := json.Marshal(ids)
	if err != nil {
		return err
	}
	r.EventIDs = string(data)
	return nil
}

// GetSubscriptionIDs 获取关联订阅源ID列表
func (r *Report) GetSubscriptionIDs() []uint {
	if r.SubscriptionIDs == "" {
		return nil
	}
	var ids []uint
	json.Unmarshal([]byte(r.SubscriptionIDs), &ids)
	return ids
}

// SetSubscriptionIDs 设置关联订阅源ID列表
func (r *Report) SetSubscriptionIDs(ids []uint) error {
	data, err := json.Marshal(ids)
	if err != nil {
		return err
	}
	r.SubscriptionIDs = string(data)
	return nil
}

// IsCompleted 检查报告是否已完成
func (r *Report) IsCompleted() bool {
	return r.Status == ReportStatusCompleted
}

// ReportTemplate 报告模板模型
type ReportTemplate struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Description string         `gorm:"size:500" json:"description"`
	Type        ReportType     `gorm:"size:20;index" json:"type"`
	Content     string         `gorm:"type:longtext" json:"content"`
	IsDefault   bool           `gorm:"default:false" json:"is_default"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// TableName 指定表名
func (ReportTemplate) TableName() string {
	return "report_templates"
}
