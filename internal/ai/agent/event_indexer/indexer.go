package event_indexer

import (
	"SuperBizAgent/internal/ai/indexer"
	"SuperBizAgent/internal/model"
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/eino/schema"
)

// EventIndexer 将处理后的安全事件索引到 Milvus 向量库
type EventIndexer struct {
	ctx context.Context
}

func NewEventIndexer(ctx context.Context) *EventIndexer {
	return &EventIndexer{ctx: ctx}
}

// IndexEvent 将安全事件和AI分析结论索引到Milvus
func (ei *EventIndexer) IndexEvent(event *model.SecurityEvent, analysis string) error {
	idx, err := indexer.NewMilvusIndexer(ei.ctx)
	if err != nil {
		return fmt.Errorf("init indexer: %w", err)
	}

	content := buildContent(event, analysis)

	doc := &schema.Document{
		ID:      fmt.Sprintf("event-%d", event.ID),
		Content: content,
		MetaData: map[string]interface{}{
			"type":       "security_event",
			"event_id":   event.ID,
			"severity":   string(event.Severity),
			"cve_id":     event.CVEID,
			"risk_score": event.RiskScore,
			"indexed_at": time.Now().Format(time.RFC3339),
		},
	}

	_, err = idx.Store(ei.ctx, []*schema.Document{doc})
	return err
}

// buildContent 构建可搜索的文本内容，截断至7000字符
func buildContent(event *model.SecurityEvent, analysis string) string {
	parts := fmt.Sprintf("标题: %s\nCVE: %s\n严重程度: %s\n厂商: %s\n产品: %s\n",
		event.Title, event.CVEID, string(event.Severity),
		event.AffectedVendor, event.AffectedProduct)

	desc := event.Description
	if len(desc) > 3000 {
		desc = desc[:3000] + "..."
	}
	parts += fmt.Sprintf("描述: %s\n", desc)

	if analysis != "" {
		remaining := 7000 - len(parts)
		if remaining > 0 {
			if len(analysis) > remaining {
				analysis = analysis[:remaining] + "..."
			}
			parts += fmt.Sprintf("AI分析: %s", analysis)
		}
	}

	if len(parts) > 7000 {
		parts = parts[:7000]
	}
	return parts
}
