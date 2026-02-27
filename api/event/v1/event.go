package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ListReq 事件列表请求
type ListReq struct {
	g.Meta   `path:"/event" method:"get" tags:"Event" summary:"获取事件列表"`
	Page     int    `json:"page" d:"1" dc:"页码"`
	PageSize int    `json:"page_size" d:"20" dc:"每页数量"`
	Severity string `json:"severity" dc:"严重程度"`
	Status   string `json:"status" dc:"事件状态"`
	Keyword  string `json:"keyword" dc:"关键词"`
}

// ListRes 事件列表响应
type ListRes struct {
	List  []EventItem `json:"list"`
	Total int64       `json:"total"`
}

// EventItem 事件列表项
type EventItem struct {
	ID        uint    `json:"id"`
	Title     string  `json:"title"`
	Severity  string  `json:"severity"`
	Status    string  `json:"status"`
	Source    string  `json:"source"`
	SourceURL string  `json:"source_url"`
	CVEID     string  `json:"cve_id"`
	CVSSScore float64 `json:"cvss_score"`
	EventTime string  `json:"event_time"`
	RiskScore int     `json:"risk_score"`
	IsStarred bool    `json:"is_starred"`
}

// GetReq 获取事件详情请求
type GetReq struct {
	g.Meta `path:"/event/{id}" method:"get" tags:"Event" summary:"获取事件详情"`
	ID     uint `json:"id" in:"path" v:"required" dc:"事件ID"`
}

// GetRes 获取事件详情响应
type GetRes struct {
	ID             uint    `json:"id"`
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	Severity       string  `json:"severity"`
	Status         string  `json:"status"`
	CVEID          string  `json:"cve_id"`
	CVSSScore      float64 `json:"cvss_score"`
	SourceURL      string  `json:"source_url"`
	EventTime      string  `json:"event_time"`
	CreatedAt      string  `json:"created_at"`
	RiskScore      int     `json:"risk_score"`
	Recommendation string  `json:"recommendation"`
	AffectedAssets int     `json:"affected_assets"`
}

// UpdateStatusReq 更新事件状态请求
type UpdateStatusReq struct {
	g.Meta `path:"/event/{id}/status" method:"put" tags:"Event" summary:"更新事件状态"`
	ID     uint   `json:"id" in:"path" v:"required" dc:"事件ID"`
	Status string `json:"status" v:"required" dc:"事件状态"`
}

// UpdateStatusRes 更新事件状态响应
type UpdateStatusRes struct{}

// BatchUpdateStatusReq 批量更新事件状态请求
type BatchUpdateStatusReq struct {
	g.Meta `path:"/event/batch/status" method:"post" tags:"Event" summary:"批量更新事件状态"`
	IDs    []uint `json:"ids" v:"required" dc:"事件ID列表"`
	Status string `json:"status" v:"required" dc:"事件状态"`
}

// BatchUpdateStatusRes 批量更新事件状态响应
type BatchUpdateStatusRes struct{}

// StatsReq 获取事件统计请求
type StatsReq struct {
	g.Meta `path:"/event/stats" method:"get" tags:"Event" summary:"获取事件统计"`
}

// StatsRes 获取事件统计响应
type StatsRes struct {
	Total        int64            `json:"total"`
	TodayCount   int64            `json:"today_count"`
	CriticalCount int64           `json:"critical_count"`
	HighCount    int64            `json:"high_count"`
	BySeverity   map[string]int64 `json:"by_severity"`
	ByStatus     map[string]int64 `json:"by_status"`
}

// TrendReq 获取事件趋势请求
type TrendReq struct {
	g.Meta `path:"/event/trend" method:"get" tags:"Event" summary:"获取事件趋势"`
	Days   int `json:"days" d:"7" dc:"天数"`
}

// TrendRes 获取事件趋势响应
type TrendRes struct {
	List []TrendItem `json:"list"`
}

// TrendItem 趋势项
type TrendItem struct {
	Date     string `json:"date"`
	Total    int64  `json:"total"`
	Critical int64  `json:"critical"`
	High     int64  `json:"high"`
	Medium   int64  `json:"medium"`
	Low      int64  `json:"low"`
	Info     int64  `json:"info"`
}

// DeleteAllReq 删除所有事件请求
type DeleteAllReq struct {
	g.Meta `path:"/event/all" method:"delete" tags:"Event" summary:"删除所有事件"`
}

// DeleteAllRes 删除所有事件响应
type DeleteAllRes struct{}

// AnalyzeReq AI分析事件请求
type AnalyzeReq struct {
	g.Meta `path:"/event/{id}/analyze" method:"post" tags:"Event" summary:"AI分析事件"`
	ID     uint `json:"id" in:"path" v:"required"`
}

// AnalyzeRes AI分析事件响应
type AnalyzeRes struct {
	RiskScore      int    `json:"risk_score"`
	Severity       string `json:"severity"`
	Recommendation string `json:"recommendation"`
}

// ProcessPipelineReq 事件处理流水线请求
type ProcessPipelineReq struct {
	g.Meta `path:"/event/pipeline/process" method:"post" tags:"Event" summary:"多Agent流水线处理事件"`
}

// ProcessPipelineRes 事件处理流水线响应
type ProcessPipelineRes struct {
	TotalCount  int            `json:"total_count"`
	DedupCount  int            `json:"dedup_count"`
	NewCount    int            `json:"new_count"`
	ProcessedAt string         `json:"processed_at"`
	Steps       []PipelineStep `json:"steps"`
}

// PipelineStep Agent处理步骤
type PipelineStep struct {
	Agent   string `json:"agent"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Count   int    `json:"count"`
}
