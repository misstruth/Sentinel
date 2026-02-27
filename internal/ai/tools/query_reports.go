package tools

import (
	"context"
	"encoding/json"

	"SuperBizAgent/internal/database"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

type ReportQueryTool struct{}

func NewReportQueryTool() tool.BaseTool {
	return &ReportQueryTool{}
}

func (t *ReportQueryTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "query_reports",
		Desc: `查询安全报告。返回报告标题、类型、状态和创建时间。

使用场景:
- 用户问"最近的报告" → 不传参数
- 用户问"周报/月报" → type="weekly"或"monthly"

返回: {count, items: [{id, title, type, status, created_at}]}`,
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"type":  {Type: schema.String, Desc: "报告类型: daily(日报), weekly(周报), monthly(月报)"},
			"limit": {Type: schema.Integer, Desc: "返回数量，1-10，默认5"},
		}),
	}, nil
}

func (t *ReportQueryTool) InvokableRun(ctx context.Context, args string, opts ...tool.Option) (string, error) {
	var params struct {
		Type  string `json:"type"`
		Limit int    `json:"limit"`
	}
	json.Unmarshal([]byte(args), &params)
	if params.Limit <= 0 {
		params.Limit = 5
	}

	db := database.GetDB()
	query := db.Table("reports").Select("id, title, type, status, created_at").Order("created_at DESC").Limit(params.Limit)
	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}

	var results []map[string]interface{}
	query.Find(&results)

	data, _ := json.Marshal(map[string]interface{}{"success": true, "count": len(results), "items": results})
	return string(data), nil
}
