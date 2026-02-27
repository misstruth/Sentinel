package report

import (
	"context"

	v1 "SuperBizAgent/api/report/v1"
)

// Get 获取报告
func (c *Controller) Get(ctx context.Context, req *v1.GetReq) (*v1.GetRes, error) {
	report, err := c.service.Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &v1.GetRes{
		ID:            report.ID,
		Title:         report.Title,
		Type:          string(report.Type),
		Status:        string(report.Status),
		Content:       report.Content,
		Summary:       report.Summary,
		EventCount:    report.EventCount,
		CriticalCount: report.CriticalCount,
		HighCount:     report.HighCount,
		StartTime:     report.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:       report.EndTime.Format("2006-01-02 15:04:05"),
		GeneratedBy:   report.GeneratedBy,
		ErrorMsg:      report.ErrorMsg,
		CreatedAt:     report.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
