package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// GenerateReq 生成报告请求
type GenerateReq struct {
	g.Meta    `path:"/report/generate" method:"post" tags:"Report" summary:"生成报告"`
	Type      string `json:"type" v:"required" dc:"报告类型: daily/weekly/monthly"`
	Title     string `json:"title" v:"required" dc:"报告标题"`
	StartTime string `json:"start_time" v:"required" dc:"开始时间"`
	EndTime   string `json:"end_time" v:"required" dc:"结束时间"`
}

// GenerateRes 生成报告响应
type GenerateRes struct {
	ReportID uint   `json:"report_id"`
	Title    string `json:"title"`
	Summary  string `json:"summary"`
	Status   string `json:"status"`
}

// GetReq 获取报告请求
type GetReq struct {
	g.Meta `path:"/report/:id" method:"get" tags:"Report" summary:"获取报告"`
	ID     uint `json:"id" v:"required" dc:"报告ID"`
}

// GetRes 获取报告响应
type GetRes struct {
	ID            uint   `json:"id"`
	Title         string `json:"title"`
	Type          string `json:"type"`
	Status        string `json:"status"`
	Content       string `json:"content"`
	Summary       string `json:"summary"`
	EventCount    int    `json:"event_count"`
	CriticalCount int    `json:"critical_count"`
	HighCount     int    `json:"high_count"`
	StartTime     string `json:"start_time"`
	EndTime       string `json:"end_time"`
	GeneratedBy   string `json:"generated_by"`
	ErrorMsg      string `json:"error_msg"`
	CreatedAt     string `json:"created_at"`
}

// ListReq 获取报告列表请求
type ListReq struct {
	g.Meta   `path:"/report" method:"get" tags:"Report" summary:"获取报告列表"`
	Page     int    `json:"page" d:"1" dc:"页码"`
	PageSize int    `json:"page_size" d:"10" dc:"每页数量"`
	Type     string `json:"type" dc:"报告类型"`
}

// ListRes 获取报告列表响应
type ListRes struct {
	List  []ReportItem `json:"list"`
	Total int64        `json:"total"`
}

// ReportItem 报告列表项
type ReportItem struct {
	ID            uint   `json:"id"`
	Title         string `json:"title"`
	Type          string `json:"type"`
	Status        string `json:"status"`
	Summary       string `json:"summary"`
	EventCount    int    `json:"event_count"`
	CriticalCount int    `json:"critical_count"`
	HighCount     int    `json:"high_count"`
	GeneratedBy   string `json:"generated_by"`
	ErrorMsg      string `json:"error_msg"`
	CreatedAt     string `json:"created_at"`
}

// DeleteReq 删除报告请求
type DeleteReq struct {
	g.Meta `path:"/report/:id" method:"delete" tags:"Report" summary:"删除报告"`
	ID     uint `json:"id" v:"required" dc:"报告ID"`
}

// DeleteRes 删除报告响应
type DeleteRes struct {
	Success bool `json:"success"`
}

// TemplateCreateReq 创建模板请求
type TemplateCreateReq struct {
	g.Meta      `path:"/report/template" method:"post" tags:"ReportTemplate" summary:"创建模板"`
	Name        string `json:"name" v:"required" dc:"模板名称"`
	Description string `json:"description" dc:"模板描述"`
	Type        string `json:"type" v:"required" dc:"报告类型"`
	Content     string `json:"content" v:"required" dc:"模板内容"`
	IsDefault   bool   `json:"is_default" dc:"是否默认"`
}

// TemplateCreateRes 创建模板响应
type TemplateCreateRes struct {
	ID uint `json:"id"`
}

// TemplateListReq 获取模板列表请求
type TemplateListReq struct {
	g.Meta `path:"/report/template" method:"get" tags:"ReportTemplate" summary:"获取模板列表"`
	Type   string `json:"type" dc:"报告类型"`
}

// TemplateListRes 获取模板列表响应
type TemplateListRes struct {
	List []TemplateItem `json:"list"`
}

// TemplateItem 模板列表项
type TemplateItem struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	IsDefault   bool   `json:"is_default"`
}

// ExportReq 导出报告请求
type ExportReq struct {
	g.Meta `path:"/report/:id/export" method:"get" tags:"Report" summary:"导出报告"`
	ID     uint   `json:"id" v:"required" dc:"报告ID"`
	Format string `json:"format" d:"markdown" dc:"导出格式: markdown/html/json"`
}

// ExportRes 导出报告响应
type ExportRes struct {
	Content  string `json:"content"`
	Filename string `json:"filename"`
}

// TemplateGetReq 获取模板详情请求
type TemplateGetReq struct {
	g.Meta `path:"/report/template/{id}" method:"get" tags:"ReportTemplate" summary:"获取模板详情"`
	ID     uint `json:"id" in:"path" v:"required" dc:"模板ID"`
}

// TemplateGetRes 获取模板详情响应
type TemplateGetRes struct {
	*TemplateItem
	Content string `json:"content"`
}

// TemplateUpdateReq 更新模板请求
type TemplateUpdateReq struct {
	g.Meta      `path:"/report/template/{id}" method:"put" tags:"ReportTemplate" summary:"更新模板"`
	ID          uint   `json:"id" in:"path" v:"required" dc:"模板ID"`
	Name        string `json:"name" dc:"模板名称"`
	Description string `json:"description" dc:"模板描述"`
	Content     string `json:"content" dc:"模板内容"`
	IsDefault   bool   `json:"is_default" dc:"是否默认"`
}

// TemplateUpdateRes 更新模板响应
type TemplateUpdateRes struct{}

// TemplateDeleteReq 删除模板请求
type TemplateDeleteReq struct {
	g.Meta `path:"/report/template/{id}" method:"delete" tags:"ReportTemplate" summary:"删除模板"`
	ID     uint `json:"id" in:"path" v:"required" dc:"模板ID"`
}

// TemplateDeleteRes 删除模板响应
type TemplateDeleteRes struct{}
