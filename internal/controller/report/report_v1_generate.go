package report

import (
	"context"
	"time"

	v1 "SuperBizAgent/api/report/v1"
	"SuperBizAgent/internal/ai/agent/report_generator"
	"SuperBizAgent/internal/model"
)

// Generate 生成报告
func (c *Controller) Generate(ctx context.Context, req *v1.GenerateReq) (*v1.GenerateRes, error) {
	startTime, _ := time.Parse("2006-01-02 15:04:05", req.StartTime)
	endTime, _ := time.Parse("2006-01-02 15:04:05", req.EndTime)

	genReq := &report_generator.ReportRequest{
		Type:        model.ReportType(req.Type),
		Title:       req.Title,
		StartTime:   startTime,
		EndTime:     endTime,
		GeneratedBy: "api",
	}

	resp, err := c.service.Generate(ctx, genReq)
	if err != nil {
		return nil, err
	}

	return &v1.GenerateRes{
		ReportID: resp.ReportID,
		Title:    resp.Title,
		Summary:  resp.Summary,
		Status:   resp.Status,
	}, nil
}
