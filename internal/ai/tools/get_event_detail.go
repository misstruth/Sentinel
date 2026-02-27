package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"SuperBizAgent/internal/database"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

type EventDetailTool struct{}

func NewEventDetailTool() tool.BaseTool {
	return &EventDetailTool{}
}

func (t *EventDetailTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "get_event_detail",
		Desc: `获取安全事件的完整详情，包括描述、原文链接、CVE信息等。

使用场景:
- 用户想了解某个事件的详细内容
- 用户想查看事件的原文链接
- 用户问"事件186是什么" → event_id=186

返回格式: {id, title, description, severity, status, source, source_url, cve_id, cvss_score, event_time, recommendation}`,
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"event_id": {Type: schema.Integer, Desc: "事件ID，必填", Required: true},
		}),
	}, nil
}

func (t *EventDetailTool) InvokableRun(ctx context.Context, args string, opts ...tool.Option) (string, error) {
	var params struct {
		EventID int `json:"event_id"`
	}
	json.Unmarshal([]byte(args), &params)

	fmt.Printf("[get_event_detail] called with args: %s, event_id: %d\n", args, params.EventID)

	if params.EventID <= 0 {
		return `{"error": "event_id is required"}`, nil
	}

	db := database.GetDB()
	var event struct {
		ID             uint    `json:"id"`
		Title          string  `json:"title"`
		Description    string  `json:"description"`
		Severity       string  `json:"severity"`
		Status         string  `json:"status"`
		Source         string  `json:"source"`
		SourceURL      string  `json:"source_url" gorm:"column:source_url"`
		CVEID          string  `json:"cve_id" gorm:"column:cve_id"`
		CVSSScore      float64 `json:"cvss_score" gorm:"column:cvss_score"`
		EventTime      string  `json:"event_time" gorm:"column:event_time"`
		Recommendation string  `json:"recommendation"`
		RiskScore      int     `json:"risk_score" gorm:"column:risk_score"`
	}

	err := db.Table("security_events").
		Where("id = ?", params.EventID).
		First(&event).Error

	if err != nil {
		return `{"error": "event not found: ` + err.Error() + `"}`, nil
	}

	data, _ := json.Marshal(event)
	return string(data), nil
}
