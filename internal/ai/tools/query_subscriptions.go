package tools

import (
	"context"
	"encoding/json"

	"SuperBizAgent/internal/database"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

type SubscriptionQueryTool struct{}

func NewSubscriptionQueryTool() tool.BaseTool {
	return &SubscriptionQueryTool{}
}

func (t *SubscriptionQueryTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "query_subscriptions",
		Desc: `查询订阅源配置。返回订阅名称、类型、URL、状态和事件数量。

使用场景:
- 用户问"有哪些订阅源" → 不传参数
- 用户问"哪些订阅在运行" → status="active"
- 用户问"订阅配置" → 不传参数

返回: {count, items: [{id, name, source_type, source_url, status, total_events}]}`,
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"status": {
				Type: schema.String,
				Desc: "状态过滤。可选值: active(运行中), paused(暂停), disabled(禁用)。不传返回全部",
			},
		}),
	}, nil
}

func (t *SubscriptionQueryTool) InvokableRun(ctx context.Context, args string, opts ...tool.Option) (string, error) {
	var params struct {
		Status string `json:"status"`
	}
	json.Unmarshal([]byte(args), &params)

	db := database.GetDB()
	query := db.Table("subscriptions").Select("id, name, source_type, source_url, status, total_events, created_at")
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	var results []map[string]interface{}
	if err := query.Find(&results).Error; err != nil {
		return `{"error": "` + err.Error() + `"}`, nil
	}

	data, _ := json.Marshal(map[string]interface{}{
		"success": true,
		"count":   len(results),
		"items":   results,
	})
	return string(data), nil
}
