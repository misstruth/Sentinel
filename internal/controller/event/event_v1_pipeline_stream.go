package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"SuperBizAgent/internal/ai/agent/event_analysis_pipeline"
	"SuperBizAgent/internal/ai/agent/event_dedup"
	"SuperBizAgent/internal/ai/agent/event_extraction"
	"SuperBizAgent/internal/ai/agent/event_indexer"
	"SuperBizAgent/internal/ai/agent/risk_assessment"
	"SuperBizAgent/internal/database"
	"SuperBizAgent/internal/model"

	"github.com/gogf/gf/v2/net/ghttp"
)

// AgentStreamEvent SSE事件
type AgentStreamEvent struct {
	Type      string      `json:"type"`
	Agent     string      `json:"agent"`
	Status    string      `json:"status"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp string      `json:"timestamp"`
}

// PipelineStream SSE流式处理 — 真实AI多Agent协作
func (c *Controller) PipelineStream(r *ghttp.Request) {
	r.Response.Header().Set("Content-Type", "text/event-stream")
	r.Response.Header().Set("Cache-Control", "no-cache")
	r.Response.Header().Set("Connection", "keep-alive")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
	defer cancel()

	send := func(e AgentStreamEvent) {
		e.Timestamp = time.Now().Format(time.RFC3339)
		data, _ := json.Marshal(e)
		r.Response.Writef("data: %s\n\n", data)
		r.Response.Flush()
	}

	// ========== Step 1: 数据采集 ==========
	send(AgentStreamEvent{Type: "agent_start", Agent: "数据采集Agent", Status: "running", Message: "扫描待处理事件..."})

	db := database.GetDB()
	var events []model.SecurityEvent
	db.Where("status = ?", "new").Limit(10).Find(&events)

	send(AgentStreamEvent{
		Type: "agent_complete", Agent: "数据采集Agent", Status: "success",
		Message: fmt.Sprintf("发现 %d 个事件", len(events)),
		Data:    map[string]interface{}{"count": len(events), "sources": []string{"NVD", "CISA KEV", "安全客"}},
	})

	if len(events) == 0 {
		send(AgentStreamEvent{Type: "pipeline_done", Agent: "Pipeline", Status: "success", Message: "无待处理事件"})
		return
	}

	// ========== Step 2: 智能提取 ==========
	allExtracted, eventDetails := stepExtraction(ctx, events, send)

	// ========== Step 3: 去重过滤 ==========
	uniqueEvents, dedupCount := stepDedup(ctx, allExtracted, send)

	// ========== Step 4: 风险评估 ==========
	assessResults, riskData := stepRiskAssessment(ctx, events, uniqueEvents, send)

	// ========== Step 5: 解决方案Agent ==========
	stepSolutionAgent(ctx, events, assessResults, eventDetails, send)

	// 将事件详情（含解决方案）合并到riskData
	riskData["events"] = eventDetails
	_ = dedupCount

	// 发送最终风险评估汇总（前端依赖此事件获取riskData）
	send(AgentStreamEvent{
		Type: "agent_complete", Agent: "风险评估Agent", Status: "success",
		Message: fmt.Sprintf("完成%d个评估", len(events)),
		Data:    riskData,
	})

	send(AgentStreamEvent{Type: "pipeline_done", Agent: "Pipeline", Status: "success", Message: "处理完成"})
}

// stepExtraction 使用AI提取结构化事件信息
func stepExtraction(ctx context.Context, events []model.SecurityEvent, send func(AgentStreamEvent)) ([]*event_extraction.ExtractedEvent, []map[string]interface{}) {
	send(AgentStreamEvent{Type: "agent_start", Agent: "提取Agent", Status: "running", Message: "提取结构化信息..."})

	extractor := event_extraction.NewExtractor(ctx)
	var allExtracted []*event_extraction.ExtractedEvent
	eventDetails := make([]map[string]interface{}, 0, len(events))

	for i, ev := range events {
		title := ev.Title
		if len(title) > 30 {
			title = title[:30]
		}
		send(AgentStreamEvent{
			Type: "agent_progress", Agent: "提取Agent", Status: "running",
			Message: fmt.Sprintf("(%d/%d) %s", i+1, len(events), title),
		})

		raw := &event_extraction.RawEventInput{
			Source:     ev.Source,
			RawContent: ev.Title + "\n" + ev.Description,
			SourceURL:  ev.SourceURL,
		}
		result, err := extractor.Extract(raw)

		if err != nil || !result.Success {
			log.Printf("[提取Agent] event %d extract failed: %v", ev.ID, err)
			// 降级：使用原始数据构造ExtractedEvent
			allExtracted = append(allExtracted, &event_extraction.ExtractedEvent{
				Title:       ev.Title,
				Description: ev.Description,
				Severity:    string(ev.Severity),
				CVEIDs:      []string{ev.CVEID},
				Source:      ev.Source,
				SourceURL:   ev.SourceURL,
			})
		} else if len(result.Events) > 0 {
			allExtracted = append(allExtracted, result.Events...)
		}

		// 构建事件详情
		desc := ev.Description
		if len(desc) > 200 {
			desc = desc[:200] + "..."
		}
		eventDetails = append(eventDetails, map[string]interface{}{
			"id": ev.ID, "title": ev.Title, "desc": desc,
			"cve_id": ev.CVEID, "cvss": ev.CVSSScore,
			"severity": ev.Severity, "vendor": ev.AffectedVendor,
			"product": ev.AffectedProduct, "source_url": ev.SourceURL,
		})
	}

	send(AgentStreamEvent{
		Type: "agent_complete", Agent: "提取Agent", Status: "success",
		Message: fmt.Sprintf("提取完成，共%d条结构化事件", len(allExtracted)),
	})
	return allExtracted, eventDetails
}

// stepDedup 语义去重
func stepDedup(ctx context.Context, allExtracted []*event_extraction.ExtractedEvent, send func(AgentStreamEvent)) ([]*event_extraction.ExtractedEvent, int) {
	send(AgentStreamEvent{Type: "agent_start", Agent: "去重Agent", Status: "running", Message: "语义去重..."})

	dedup := event_dedup.NewDeduplicator()
	result, err := dedup.Dedup(ctx, &event_dedup.DedupInput{Events: allExtracted})
	if err != nil {
		log.Printf("[去重Agent] dedup failed: %v", err)
		send(AgentStreamEvent{Type: "agent_complete", Agent: "去重Agent", Status: "success", Message: "去重完成（降级）"})
		return allExtracted, 0
	}

	msg := fmt.Sprintf("去重完成，去除%d条重复，保留%d条", result.DupCount, len(result.UniqueEvents))
	send(AgentStreamEvent{Type: "agent_complete", Agent: "去重Agent", Status: "success", Message: msg})
	return result.UniqueEvents, result.DupCount
}

// stepRiskAssessment 使用AI进行风险评估
func stepRiskAssessment(
	ctx context.Context,
	events []model.SecurityEvent,
	uniqueEvents []*event_extraction.ExtractedEvent,
	send func(AgentStreamEvent),
) (map[uint]*risk_assessment.AssessmentResult, map[string]interface{}) {
	send(AgentStreamEvent{Type: "agent_start", Agent: "风险评估Agent", Status: "running", Message: "AI风险评估中..."})

	assessor := risk_assessment.NewAssessor(ctx)
	results := make(map[uint]*risk_assessment.AssessmentResult)

	var totalScore float64
	var maxCVSS float64
	var criticalCount, highCount int

	for i, ev := range events {
		title := ev.Title
		if len(title) > 30 {
			title = title[:30]
		}

		// 找到对应的提取事件
		var extracted *event_extraction.ExtractedEvent
		if i < len(uniqueEvents) {
			extracted = uniqueEvents[i]
		} else {
			extracted = &event_extraction.ExtractedEvent{
				Title:       ev.Title,
				Description: ev.Description,
				Severity:    string(ev.Severity),
				CVEIDs:      []string{ev.CVEID},
			}
		}

		assess, err := assessor.Assess(&risk_assessment.AssessmentInput{Event: extracted})

		if err != nil {
			log.Printf("[风险评估Agent] event %d assess failed: %v", ev.ID, err)
			assess = &risk_assessment.AssessmentResult{
				RiskScore: 50, Severity: "medium",
				Recommendation: "评估失败，需人工复核",
			}
		}

		results[ev.ID] = assess
		totalScore += float64(assess.RiskScore)
		if ev.CVSSScore > maxCVSS {
			maxCVSS = ev.CVSSScore
		}
		if assess.Severity == "critical" || ev.Severity == "critical" {
			criticalCount++
		} else if assess.Severity == "high" || ev.Severity == "high" {
			highCount++
		}

		send(AgentStreamEvent{
			Type: "agent_progress", Agent: "风险评估Agent", Status: "running",
			Message: fmt.Sprintf("(%d/%d) %s → 风险分:%d", i+1, len(events), title, assess.RiskScore),
			Data: map[string]interface{}{
				"cvss": ev.CVSSScore, "risk": assess.RiskScore,
				"cve": ev.CVEID, "severity": assess.Severity,
				"recommendation": assess.Recommendation,
				"factors":        assess.Factors,
			},
		})
	}

	avgScore := 0
	if len(events) > 0 {
		avgScore = int(totalScore / float64(len(events)))
	}

	riskData := map[string]interface{}{
		"count": len(events), "avgRisk": avgScore,
		"maxCVSS": maxCVSS, "critical": criticalCount,
		"highRisk": highCount,
	}

	return results, riskData
}

// stepSolutionAgent 使用ReAct Agent查询解决方案和历史事件
func stepSolutionAgent(
	ctx context.Context,
	events []model.SecurityEvent,
	assessResults map[uint]*risk_assessment.AssessmentResult,
	eventDetails []map[string]interface{},
	send func(AgentStreamEvent),
) {
	// 筛选高危/严重事件
	var highRiskEvents []model.SecurityEvent
	for _, ev := range events {
		assess := assessResults[ev.ID]
		if assess != nil && (assess.Severity == "critical" || assess.Severity == "high") {
			highRiskEvents = append(highRiskEvents, ev)
		} else if ev.Severity == "critical" || ev.Severity == "high" {
			highRiskEvents = append(highRiskEvents, ev)
		}
	}

	if len(highRiskEvents) == 0 {
		// 无高危事件，跳过解决方案Agent
		send(AgentStreamEvent{
			Type: "agent_start", Agent: "解决方案Agent", Status: "running",
			Message: "无高危事件，跳过深度分析",
		})
		send(AgentStreamEvent{
			Type: "agent_complete", Agent: "解决方案Agent", Status: "success",
			Message: "所有事件风险可控",
		})
		return
	}

	send(AgentStreamEvent{
		Type: "agent_start", Agent: "解决方案Agent", Status: "running",
		Message: fmt.Sprintf("对%d个高危事件进行深度分析...", len(highRiskEvents)),
	})

	// 构建ReAct Agent
	agent, err := event_analysis_pipeline.BuildEventAnalysisAgent(ctx)
	if err != nil {
		log.Printf("[解决方案Agent] build agent failed: %v", err)
		send(AgentStreamEvent{
			Type: "agent_complete", Agent: "解决方案Agent", Status: "success",
			Message: "Agent初始化失败，请检查AI服务配置",
		})
		return
	}

	idxer := event_indexer.NewEventIndexer(ctx)
	// 先检查向量库容量
	event_indexer.EnsureCapacity(ctx)

	for i, ev := range highRiskEvents {
		title := ev.Title
		if len(title) > 30 {
			title = title[:30]
		}
		send(AgentStreamEvent{
			Type: "agent_progress", Agent: "解决方案Agent", Status: "running",
			Message: fmt.Sprintf("(%d/%d) 分析: %s", i+1, len(highRiskEvents), title),
		})

		// 调用ReAct Agent
		query := buildSolutionQuery(ev, assessResults[ev.ID])
		agentCtx, agentCancel := context.WithTimeout(ctx, 60*time.Second)

		msg := &event_analysis_pipeline.UserMessage{
			ID:    fmt.Sprintf("pipeline-%d", ev.ID),
			Query: query,
		}
		result, err := agent.Invoke(agentCtx, msg)
		agentCancel()

		solution := ""
		if err != nil {
			log.Printf("[解决方案Agent] event %d invoke failed: %v", ev.ID, err)
			solution = "AI分析超时，建议人工排查"
		} else {
			solution = result.Content
		}

		// 更新事件详情中的recommendation
		for _, detail := range eventDetails {
			if id, ok := detail["id"].(uint); ok && id == ev.ID {
				detail["recommendation"] = solution
			}
		}

		send(AgentStreamEvent{
			Type: "agent_progress", Agent: "解决方案Agent", Status: "running",
			Message: fmt.Sprintf("(%d/%d) 完成: %s", i+1, len(highRiskEvents), title),
			Data: map[string]interface{}{
				"event_id": ev.ID,
				"solution": truncate(solution, 500),
			},
		})

		// 异步索引到Milvus
		go func(e model.SecurityEvent, s string) {
			if err := idxer.IndexEvent(&e, s); err != nil {
				log.Printf("[event_indexer] index event %d failed: %v", e.ID, err)
			}
		}(ev, solution)
	}

	send(AgentStreamEvent{
		Type: "agent_complete", Agent: "解决方案Agent", Status: "success",
		Message: fmt.Sprintf("完成%d个高危事件的深度分析", len(highRiskEvents)),
	})
}

// buildSolutionQuery 构建解决方案查询
func buildSolutionQuery(ev model.SecurityEvent, assess *risk_assessment.AssessmentResult) string {
	var parts []string
	parts = append(parts, fmt.Sprintf("分析安全事件并给出解决方案: %s", ev.Title))
	if ev.CVEID != "" {
		parts = append(parts, fmt.Sprintf("CVE编号: %s", ev.CVEID))
	}
	if assess != nil && assess.Recommendation != "" {
		parts = append(parts, fmt.Sprintf("初步建议: %s", assess.Recommendation))
	}
	desc := ev.Description
	if len(desc) > 300 {
		desc = desc[:300] + "..."
	}
	if desc != "" {
		parts = append(parts, fmt.Sprintf("描述: %s", desc))
	}
	parts = append(parts, "请搜索历史相似事件和内部知识库，给出具体的处置步骤和修复方案。")
	return strings.Join(parts, "\n")
}

// truncate 截断字符串
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
