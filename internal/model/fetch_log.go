package model

import (
	"time"
)

// FetchStatus 抓取状态
type FetchStatus string

const (
	FetchStatusSuccess FetchStatus = "success" // 成功
	FetchStatusFailed  FetchStatus = "failed"  // 失败
	FetchStatusTimeout FetchStatus = "timeout" // 超时
)

// FetchLog 抓取日志模型
type FetchLog struct {
	ID             uint        `gorm:"primaryKey" json:"id"`
	SubscriptionID uint        `gorm:"index" json:"subscription_id"`
	Status         FetchStatus `gorm:"size:20" json:"status"`
	EventCount     int         `gorm:"default:0" json:"event_count"`
	ErrorMsg       string      `gorm:"type:text" json:"error_msg"`
	Duration       int         `gorm:"default:0" json:"duration"` // 耗时(毫秒)
	CreatedAt      time.Time   `json:"created_at"`
}

// TableName 指定表名
func (FetchLog) TableName() string {
	return "fetch_logs"
}
