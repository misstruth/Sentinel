package report_generator

import (
	"SuperBizAgent/internal/database"
	"SuperBizAgent/internal/model"
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/cloudwego/eino/schema"
)

// ReportGenerator 报告生成器
type ReportGenerator struct {
	ctx context.Context
}

// NewReportGenerator 创建报告生成器
func NewReportGenerator(ctx context.Context) *ReportGenerator {
	return &ReportGenerator{ctx: ctx}
}

// Generate 生成报告
func (g *ReportGenerator) Generate(req *ReportRequest) (*ReportResponse, error) {
	// 创建报告记录
	report := &model.Report{
		Title:       req.Title,
		Type:        req.Type,
		Status:      model.ReportStatusGenerating,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		TemplateID:  req.TemplateID,
		GeneratedBy: req.GeneratedBy,
	}

	db := database.GetDB()
	if err := db.Create(report).Error; err != nil {
		return nil, fmt.Errorf("创建报告记录失败: %w", err)
	}

	// 获取事件数据
	events, err := g.fetchEvents(req)
	if err != nil {
		g.updateReportStatus(report.ID, model.ReportStatusFailed, err.Error())
		return nil, err
	}

	// 生成报告内容
	content, summary, err := g.generateContent(req, events)
	if err != nil {
		g.updateReportStatus(report.ID, model.ReportStatusFailed, err.Error())
		return nil, err
	}

	// 更新报告
	report.Content = content
	report.Summary = summary
	report.Status = model.ReportStatusCompleted
	report.EventCount = len(events)
	report.CriticalCount = g.countBySeverity(events, model.SeverityCritical)
	report.HighCount = g.countBySeverity(events, model.SeverityHigh)

	if err := db.Save(report).Error; err != nil {
		return nil, fmt.Errorf("保存报告失败: %w", err)
	}

	return &ReportResponse{
		ReportID: report.ID,
		Title:    report.Title,
		Summary:  summary,
		Content:  content,
		Status:   string(report.Status),
	}, nil
}

// fetchEvents 获取事件数据
func (g *ReportGenerator) fetchEvents(req *ReportRequest) ([]model.SecurityEvent, error) {
	db := database.GetDB()
	var events []model.SecurityEvent

	query := db.Where("event_time BETWEEN ? AND ?", req.StartTime, req.EndTime)

	if len(req.EventIDs) > 0 {
		query = query.Where("id IN ?", req.EventIDs)
	}

	if err := query.Order("severity ASC, event_time DESC").Find(&events).Error; err != nil {
		return nil, fmt.Errorf("查询事件失败: %w", err)
	}

	return events, nil
}

// generateContent 使用LLM生成报告内容
func (g *ReportGenerator) generateContent(req *ReportRequest, events []model.SecurityEvent) (string, string, error) {
	// 获取提示词模板
	promptTpl := g.getPromptTemplate(req.Type)

	// 准备模板数据
	data := map[string]interface{}{
		"EventData": g.formatEvents(events),
		"StartTime": req.StartTime.Format("2006-01-02 15:04"),
		"EndTime":   req.EndTime.Format("2006-01-02 15:04"),
	}

	// 渲染提示词
	prompt, err := g.renderTemplate(promptTpl, data)
	if err != nil {
		return "", "", fmt.Errorf("渲染提示词失败: %w", err)
	}

	// 调用LLM生成
	content, err := g.callLLM(prompt)
	if err != nil {
		return "", "", fmt.Errorf("LLM生成失败: %w", err)
	}

	// 生成摘要
	summary := g.extractSummary(content)

	return content, summary, nil
}

// callLLM 调用LLM生成内容
func (g *ReportGenerator) callLLM(prompt string) (string, error) {
	chatModel, err := newReportModel(g.ctx)
	if err != nil {
		return "", err
	}

	messages := []*schema.Message{
		schema.UserMessage(prompt),
	}

	resp, err := chatModel.Generate(g.ctx, messages)
	if err != nil {
		return "", err
	}

	return resp.Content, nil
}

// getPromptTemplate 获取提示词模板
func (g *ReportGenerator) getPromptTemplate(reportType model.ReportType) string {
	switch reportType {
	case model.ReportTypeDaily:
		return DailyReportPrompt
	case model.ReportTypeWeekly:
		return WeeklyReportPrompt
	default:
		return DailyReportPrompt
	}
}

// renderTemplate 渲染模板
func (g *ReportGenerator) renderTemplate(tpl string, data map[string]interface{}) (string, error) {
	t, err := template.New("prompt").Parse(tpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// formatEvents 格式化事件数据
func (g *ReportGenerator) formatEvents(events []model.SecurityEvent) string {
	var buf bytes.Buffer
	for i, e := range events {
		buf.WriteString(fmt.Sprintf("%d. [%s] %s\n", i+1, e.Severity, e.Title))
		buf.WriteString(fmt.Sprintf("   时间: %s\n", e.EventTime.Format("2006-01-02 15:04")))
		if e.Description != "" {
			buf.WriteString(fmt.Sprintf("   描述: %s\n", e.Description))
		}
		buf.WriteString("\n")
	}
	return buf.String()
}

// extractSummary 提取摘要
func (g *ReportGenerator) extractSummary(content string) string {
	if len(content) > 500 {
		return content[:500] + "..."
	}
	return content
}

// countBySeverity 按严重程度统计事件数
func (g *ReportGenerator) countBySeverity(events []model.SecurityEvent, severity model.SeverityLevel) int {
	count := 0
	for _, e := range events {
		if e.Severity == severity {
			count++
		}
	}
	return count
}

// updateReportStatus 更新报告状态
func (g *ReportGenerator) updateReportStatus(id uint, status model.ReportStatus, errMsg string) {
	db := database.GetDB()
	db.Model(&model.Report{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":    status,
		"error_msg": errMsg,
	})
}
