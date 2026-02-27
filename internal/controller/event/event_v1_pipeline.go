package event

import (
	"context"
	"fmt"
	"time"

	v1 "SuperBizAgent/api/event/v1"
	"SuperBizAgent/internal/ai/agent/event_extraction"
	"SuperBizAgent/internal/ai/agent/event_pipeline"
	"SuperBizAgent/internal/database"
	"SuperBizAgent/internal/model"
)

// ProcessPipeline 多Agent流水线处理事件
func (c *Controller) ProcessPipeline(ctx context.Context, req *v1.ProcessPipelineReq) (*v1.ProcessPipelineRes, error) {
	db := database.GetDB()
	steps := []v1.PipelineStep{}

	// Step 1: 获取未处理的事件
	var events []model.SecurityEvent
	db.Where("risk_score = 0 OR risk_score IS NULL").Limit(50).Find(&events)

	steps = append(steps, v1.PipelineStep{
		Agent:   "数据采集",
		Status:  "completed",
		Message: fmt.Sprintf("获取到 %d 个待处理事件", len(events)),
		Count:   len(events),
	})

	if len(events) == 0 {
		return &v1.ProcessPipelineRes{
			ProcessedAt: time.Now().Format(time.RFC3339),
			Steps:       steps,
		}, nil
	}

	// 转换为Pipeline输入
	var rawEvents []*event_extraction.RawEventInput
	for _, e := range events {
		rawEvents = append(rawEvents, &event_extraction.RawEventInput{
			Source:     "database",
			RawContent: e.Title + "\n" + e.Description,
			SourceURL:  e.SourceURL,
		})
	}

	// Step 2-4: 执行Pipeline
	pipeline := event_pipeline.NewPipeline(ctx)
	result, err := pipeline.Process(&event_pipeline.PipelineInput{RawEvents: rawEvents})
	if err != nil {
		return nil, err
	}

	// 添加Agent处理步骤
	steps = append(steps,
		v1.PipelineStep{Agent: "提取Agent", Status: "completed", Message: "结构化提取完成", Count: result.TotalCount},
		v1.PipelineStep{Agent: "去重Agent", Status: "completed", Message: fmt.Sprintf("去重 %d 个重复事件", result.DedupCount), Count: result.DedupCount},
		v1.PipelineStep{Agent: "风险评估Agent", Status: "completed", Message: "风险评分完成", Count: len(result.Events)},
	)

	// 更新数据库
	newCount := 0
	for i, pe := range result.Events {
		if i < len(events) {
			db.Model(&events[i]).Updates(map[string]interface{}{
				"risk_score":     pe.RiskScore,
				"recommendation": pe.Recommendation,
				"severity":       pe.Severity,
			})
			newCount++
		}
	}

	steps = append(steps, v1.PipelineStep{
		Agent:   "数据持久化",
		Status:  "completed",
		Message: fmt.Sprintf("更新 %d 个事件", newCount),
		Count:   newCount,
	})

	return &v1.ProcessPipelineRes{
		TotalCount:  result.TotalCount,
		DedupCount:  result.DedupCount,
		NewCount:    newCount,
		ProcessedAt: result.ProcessedAt,
		Steps:       steps,
	}, nil
}
