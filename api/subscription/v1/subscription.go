package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// CreateReq 创建订阅请求
type CreateReq struct {
	g.Meta      `path:"/subscriptions" method:"post" summary:"创建订阅"`
	Name        string `json:"name" v:"required" dc:"订阅名称"`
	Description string `json:"description" dc:"订阅描述"`
	SourceType  string `json:"source_type" v:"required" dc:"源类型"`
	SourceURL   string `json:"source_url" dc:"源地址"`
	CronExpr    string `json:"cron_expr" dc:"Cron表达式"`
	Config      string `json:"config" dc:"JSON配置"`
}

// CreateRes 创建订阅响应
type CreateRes struct {
	ID uint `json:"id"`
}

// GetReq 获取订阅详情请求
type GetReq struct {
	g.Meta `path:"/subscriptions/{id}" method:"get" summary:"获取订阅详情"`
	ID     uint `json:"id" in:"path" v:"required" dc:"订阅ID"`
}

// GetRes 获取订阅详情响应
type GetRes struct {
	*SubscriptionItem
}

// UpdateReq 更新订阅请求
type UpdateReq struct {
	g.Meta      `path:"/subscriptions/{id}" method:"put" summary:"更新订阅"`
	ID          uint   `json:"id" in:"path" v:"required" dc:"订阅ID"`
	Name        string `json:"name" dc:"订阅名称"`
	Description string `json:"description" dc:"订阅描述"`
	SourceURL   string `json:"source_url" dc:"源地址"`
	Status      string `json:"status" dc:"状态"`
	CronExpr    string `json:"cron_expr" dc:"Cron表达式"`
	Config      string `json:"config" dc:"JSON配置"`
}

// UpdateRes 更新订阅响应
type UpdateRes struct{}

// DeleteReq 删除订阅请求
type DeleteReq struct {
	g.Meta `path:"/subscriptions/{id}" method:"delete" summary:"删除订阅"`
	ID     uint `json:"id" in:"path" v:"required" dc:"订阅ID"`
}

// DeleteRes 删除订阅响应
type DeleteRes struct{}

// ListReq 订阅列表请求
type ListReq struct {
	g.Meta     `path:"/subscriptions" method:"get" summary:"订阅列表"`
	Page       int    `json:"page" in:"query" dc:"页码"`
	PageSize   int    `json:"page_size" in:"query" dc:"每页数量"`
	SourceType string `json:"source_type" in:"query" dc:"源类型"`
	Status     string `json:"status" in:"query" dc:"状态"`
	Keyword    string `json:"keyword" in:"query" dc:"关键词"`
}

// ListRes 订阅列表响应
type ListRes struct {
	Items []*SubscriptionItem `json:"items"`
	Total int64               `json:"total"`
	Page  int                 `json:"page"`
	Size  int                 `json:"size"`
}

// SubscriptionItem 订阅项
type SubscriptionItem struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SourceType  string `json:"source_type"`
	SourceURL   string `json:"source_url"`
	Status      string `json:"status"`
	CronExpr    string `json:"cron_expr"`
	TotalEvents int    `json:"total_events"`
	LastFetchAt string `json:"last_fetch_at"`
	CreatedAt   string `json:"created_at"`
}

// PauseReq 暂停订阅请求
type PauseReq struct {
	g.Meta `path:"/subscriptions/{id}/pause" method:"post" summary:"暂停订阅"`
	ID     uint `json:"id" in:"path" v:"required" dc:"订阅ID"`
}

// PauseRes 暂停订阅响应
type PauseRes struct{}

// ResumeReq 恢复订阅请求
type ResumeReq struct {
	g.Meta `path:"/subscriptions/{id}/resume" method:"post" summary:"恢复订阅"`
	ID     uint `json:"id" in:"path" v:"required" dc:"订阅ID"`
}

// ResumeRes 恢复订阅响应
type ResumeRes struct{}

// DisableReq 禁用订阅请求
type DisableReq struct {
	g.Meta `path:"/subscriptions/{id}/disable" method:"post" summary:"禁用订阅"`
	ID     uint `json:"id" in:"path" v:"required" dc:"订阅ID"`
}

// DisableRes 禁用订阅响应
type DisableRes struct{}

// FetchLogsReq 获取抓取日志请求
type FetchLogsReq struct {
	g.Meta   `path:"/subscriptions/{id}/logs" method:"get" summary:"获取抓取日志"`
	ID       uint `json:"id" in:"path" v:"required" dc:"订阅ID"`
	Page     int  `json:"page" in:"query" d:"1" dc:"页码"`
	PageSize int  `json:"page_size" in:"query" d:"20" dc:"每页数量"`
}

// FetchLogsRes 获取抓取日志响应
type FetchLogsRes struct {
	List  []*FetchLogItem `json:"list"`
	Total int64           `json:"total"`
	Page  int             `json:"page"`
	Size  int             `json:"size"`
}

// FetchLogItem 抓取日志项
type FetchLogItem struct {
	ID             uint   `json:"id"`
	SubscriptionID uint   `json:"subscription_id"`
	Status         string `json:"status"`
	EventCount     int    `json:"event_count"`
	ErrorMsg       string `json:"error_msg"`
	Duration       int    `json:"duration"`
	CreatedAt      string `json:"created_at"`
}

// StatsReq 获取订阅统计请求
type StatsReq struct {
	g.Meta `path:"/subscriptions/{id}/stats" method:"get" summary:"获取订阅统计"`
	ID     uint `json:"id" in:"path" v:"required" dc:"订阅ID"`
}

// StatsRes 获取订阅统计响应
type StatsRes struct {
	TotalFetches  int64 `json:"total_fetches"`
	SuccessCount  int64 `json:"success_count"`
	FailedCount   int64 `json:"failed_count"`
	TotalEvents   int   `json:"total_events"`
	AvgDurationMs int   `json:"avg_duration_ms"`
}

// FetchReq 手动触发抓取请求
type FetchReq struct {
	g.Meta `path:"/subscriptions/{id}/fetch" method:"post" summary:"手动触发抓取"`
	ID     uint `json:"id" in:"path" v:"required" dc:"订阅ID"`
}

// FetchRes 手动触发抓取响应
type FetchRes struct {
	FetchedCount int    `json:"fetched_count"`
	NewCount     int    `json:"new_count"`
	TotalEvents  int64  `json:"total_events"`
	Duration     int    `json:"duration_ms"`
	Message      string `json:"message"`
}
