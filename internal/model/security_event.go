package model

import (
	"time"

	"gorm.io/gorm"
)

// SeverityLevel 严重程度
type SeverityLevel string

const (
	SeverityCritical SeverityLevel = "critical" // 严重
	SeverityHigh     SeverityLevel = "high"     // 高危
	SeverityMedium   SeverityLevel = "medium"   // 中危
	SeverityLow      SeverityLevel = "low"      // 低危
	SeverityInfo     SeverityLevel = "info"     // 信息
)

// EventStatus 事件状态
type EventStatus string

const (
	EventStatusNew        EventStatus = "new"        // 新事件
	EventStatusProcessing EventStatus = "processing" // 处理中
	EventStatusResolved   EventStatus = "resolved"   // 已解决
	EventStatusIgnored    EventStatus = "ignored"    // 已忽略
)

// SecurityEvent 安全事件模型
type SecurityEvent struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	SubscriptionID uint           `gorm:"index" json:"subscription_id"`          // 关联订阅源
	Title          string         `gorm:"size:500;not null" json:"title"`        // 事件标题
	Description    string         `gorm:"type:text" json:"description"`          // 事件描述
	Severity       SeverityLevel  `gorm:"size:20;index" json:"severity"`         // 严重程度
	Status         EventStatus    `gorm:"size:20;default:new;index" json:"status"` // 事件状态
	SourceURL      string         `gorm:"size:500" json:"source_url"`            // 原始链接
	EventTime      time.Time      `gorm:"index" json:"event_time"`               // 事件时间
	RawData        string         `gorm:"type:text" json:"raw_data"`             // 原始数据

	// 扩展字段
	Source      string `gorm:"size:100" json:"source"`            // 来源
	IsStarred   bool   `gorm:"default:false" json:"is_starred"`   // 是否收藏
	CVEID       string `gorm:"size:50;index" json:"cve_id"`       // CVE编号
	CVSSScore   float64 `gorm:"default:0" json:"cvss_score"`      // CVSS评分
	AffectedVendor string `gorm:"size:100" json:"affected_vendor"` // 受影响厂商
	AffectedProduct string `gorm:"size:200" json:"affected_product"` // 受影响产品
	Tags        string `gorm:"size:500" json:"tags"`              // 标签(逗号分隔)
	UniqueHash  string `gorm:"size:64;uniqueIndex" json:"unique_hash"` // 去重哈希

	// AI分析字段
	RiskScore      int    `gorm:"default:0" json:"risk_score"`       // 风险评分(0-100)
	Recommendation string `gorm:"type:text" json:"recommendation"`   // 处置建议
	RelatedEvents  string `gorm:"type:text" json:"related_events"`   // 关联事件ID(JSON)

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// TableName 指定表名
func (SecurityEvent) TableName() string {
	return "security_events"
}
