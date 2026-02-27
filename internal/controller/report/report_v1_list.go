package report

import (
	"context"

	v1 "SuperBizAgent/api/report/v1"
)

// List 获取报告列表
func (c *Controller) List(ctx context.Context, req *v1.ListReq) (*v1.ListRes, error) {
	reports, total, err := c.service.List(ctx, req.Page, req.PageSize, req.Type)
	if err != nil {
		return nil, err
	}

	list := make([]v1.ReportItem, len(reports))
	for i, r := range reports {
		list[i] = v1.ReportItem{
			ID:            r.ID,
			Title:         r.Title,
			Type:          string(r.Type),
			Status:        string(r.Status),
			Summary:       r.Summary,
			EventCount:    r.EventCount,
			CriticalCount: r.CriticalCount,
			HighCount:     r.HighCount,
			GeneratedBy:   r.GeneratedBy,
			ErrorMsg:      r.ErrorMsg,
			CreatedAt:     r.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return &v1.ListRes{
		List:  list,
		Total: total,
	}, nil
}
