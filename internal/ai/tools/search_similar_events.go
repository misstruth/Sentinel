package tools

import (
	"SuperBizAgent/internal/ai/retriever"
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

type SearchSimilarEventsInput struct {
	Query string `json:"query" jsonschema:"description=搜索查询，用于在向量数据库中查找相似的历史安全事件。可以是事件标题、CVE编号、漏洞描述等"`
	Limit int    `json:"limit" jsonschema:"description=返回的最大相似事件数量，默认3"`
}

func NewSearchSimilarEventsTool() tool.InvokableTool {
	t, err := utils.InferOptionableTool(
		"search_similar_events",
		`在向量数据库中搜索历史相似安全事件。基于语义相似度匹配已处理过的安全事件。

使用场景:
- 发现新漏洞时，查找历史上类似的漏洞和处置记录
- 分析攻击模式时，查找相似的攻击事件
- 需要参考历史处置方案时使用

返回: 相似事件列表，包含标题、严重程度、CVE、相似度分数和历史处置建议`,
		func(ctx context.Context, input *SearchSimilarEventsInput, opts ...tool.Option) (string, error) {
			if input.Query == "" {
				return `{"error": "query is required"}`, nil
			}

			rr, err := retriever.NewMilvusRetriever(ctx)
			if err != nil {
				return fmt.Sprintf(`{"error": "retriever init failed: %v"}`, err), nil
			}

			docs, err := rr.Retrieve(ctx, input.Query)
			if err != nil {
				return fmt.Sprintf(`{"error": "retrieve failed: %v"}`, err), nil
			}

			limit := input.Limit
			if limit <= 0 || limit > 5 {
				limit = 3
			}

			var results []map[string]interface{}
			for i, doc := range docs {
				if i >= limit {
					break
				}
				item := map[string]interface{}{
					"content": doc.Content,
				}
				if doc.MetaData != nil {
					for k, v := range doc.MetaData {
						item[k] = v
					}
				}
				if score, ok := doc.MetaData["_score"]; ok {
					item["similarity_score"] = score
				}
				results = append(results, item)
			}

			resp := map[string]interface{}{
				"count": len(results),
				"items": results,
			}
			data, _ := json.Marshal(resp)
			return string(data), nil
		})
	if err != nil {
		return nil
	}
	return t
}
