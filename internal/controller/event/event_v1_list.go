package event

import (
	"context"

	v1 "SuperBizAgent/api/event/v1"
)

// List 获取事件列表
func (c *Controller) List(ctx context.Context, req *v1.ListReq) (*v1.ListRes, error) {
	events, total, err := c.service.List(ctx, req.Page, req.PageSize, req.Severity, req.Status, req.Keyword)
	if err != nil {
		return nil, err
	}

	list := make([]v1.EventItem, len(events))
	for i, e := range events {
		list[i] = v1.EventItem{
			ID:        e.ID,
			Title:     e.Title,
			Severity:  string(e.Severity),
			Status:    string(e.Status),
			Source:    e.Source,
			SourceURL: e.SourceURL,
			CVEID:     e.CVEID,
			CVSSScore: e.CVSSScore,
			EventTime: e.EventTime.Format("2006-01-02 15:04:05"),
			RiskScore: e.RiskScore,
			IsStarred: e.IsStarred,
		}
	}

	return &v1.ListRes{List: list, Total: total}, nil
}
