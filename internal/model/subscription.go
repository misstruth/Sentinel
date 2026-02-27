package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// SourceType 订阅源类型
type SourceType string

const (
	SourceTypeVulnerability  SourceType = "vulnerability"   // 漏洞情报
	SourceTypeThreatIntel    SourceType = "threat_intel"    // 威胁情报
	SourceTypeVendorAdvisory SourceType = "vendor_advisory" // 厂商公告
	SourceTypeAttackActivity SourceType = "attack_activity" // 攻击活动
	SourceTypeGitHubRepo     SourceType = "github_repo"     // GitHub 仓库
	SourceTypeRSS            SourceType = "rss"             // RSS 订阅
	SourceTypeWebHook        SourceType = "webhook"         // WebHook
	SourceTypeNVD            SourceType = "nvd"             // NVD 漏洞库
	SourceTypeCVE            SourceType = "cve"             // CVE 漏洞库
)

// SubscriptionStatus 订阅状态
type SubscriptionStatus string

const (
	StatusActive   SubscriptionStatus = "active"   // 活跃
	StatusPaused   SubscriptionStatus = "paused"   // 暂停
	StatusDisabled SubscriptionStatus = "disabled" // 禁用
)

// Subscription 订阅源模型
type Subscription struct {
	ID          uint               `gorm:"primaryKey" json:"id"`
	Name        string             `gorm:"size:255;not null" json:"name"`              // 订阅名称
	Description string             `gorm:"size:1000" json:"description"`               // 订阅描述
	SourceType  SourceType         `gorm:"size:50;not null;index" json:"source_type"`  // 源类型
	SourceURL   string             `gorm:"size:500" json:"source_url"`                 // 源地址
	Status      SubscriptionStatus `gorm:"size:20;default:active;index" json:"status"` // 状态
	Config      string             `gorm:"type:text" json:"config"`                    // JSON 配置

	// 调度配置
	CronExpr     string     `gorm:"size:100" json:"cron_expr"`       // Cron 表达式
	LastFetchAt  *time.Time `json:"last_fetch_at"`                   // 上次抓取时间
	NextFetchAt  *time.Time `json:"next_fetch_at"`                   // 下次抓取时间
	FetchTimeout int        `gorm:"default:30" json:"fetch_timeout"` // 抓取超时(秒)

	// 认证配置
	AuthType  string `gorm:"size:50" json:"auth_type"`  // 认证类型: none/api_key/oauth/basic
	AuthToken string `gorm:"size:500" json:"auth_token"` // 认证令牌(加密存储)

	// 过滤配置
	Keywords    string `gorm:"type:text" json:"keywords"`     // 关键词过滤(JSON数组)
	MinSeverity string `gorm:"size:20" json:"min_severity"`   // 最低严重级别
	Tags        string `gorm:"size:500" json:"tags"`          // 标签(逗号分隔)

	// 统计信息
	TotalEvents   int `gorm:"default:0" json:"total_events"`   // 总事件数
	FailedFetches int `gorm:"default:0" json:"failed_fetches"` // 失败次数

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// TableName 指定表名
func (Subscription) TableName() string {
	return "subscriptions"
}

// SubscriptionConfig 订阅源配置结构
type SubscriptionConfig struct {
	// GitHub 配置
	GitHubOwner string `json:"github_owner,omitempty"`
	GitHubRepo  string `json:"github_repo,omitempty"`
	WatchEvents []string `json:"watch_events,omitempty"` // issues, releases, security

	// RSS 配置
	RSSFeedURL string `json:"rss_feed_url,omitempty"`

	// WebHook 配置
	WebHookSecret string `json:"webhook_secret,omitempty"`

	// API 配置
	APIEndpoint string            `json:"api_endpoint,omitempty"`
	APIHeaders  map[string]string `json:"api_headers,omitempty"`

	// 通用配置
	MaxItems    int  `json:"max_items,omitempty"`
	EnableProxy bool `json:"enable_proxy,omitempty"`
	ProxyURL    string `json:"proxy_url,omitempty"`
}

// GetConfig 解析配置JSON
func (s *Subscription) GetConfig() (*SubscriptionConfig, error) {
	if s.Config == "" {
		return &SubscriptionConfig{}, nil
	}
	var config SubscriptionConfig
	err := json.Unmarshal([]byte(s.Config), &config)
	return &config, err
}

// SetConfig 设置配置JSON
func (s *Subscription) SetConfig(config *SubscriptionConfig) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	s.Config = string(data)
	return nil
}

// GetKeywords 获取关键词列表
func (s *Subscription) GetKeywords() []string {
	if s.Keywords == "" {
		return nil
	}
	var keywords []string
	json.Unmarshal([]byte(s.Keywords), &keywords)
	return keywords
}

// SetKeywords 设置关键词列表
func (s *Subscription) SetKeywords(keywords []string) error {
	data, err := json.Marshal(keywords)
	if err != nil {
		return err
	}
	s.Keywords = string(data)
	return nil
}

// IsActive 检查订阅是否活跃
func (s *Subscription) IsActive() bool {
	return s.Status == StatusActive
}
