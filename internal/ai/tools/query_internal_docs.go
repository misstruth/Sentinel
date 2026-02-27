package tools

import (
	"SuperBizAgent/internal/ai/retriever"
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

type QueryInternalDocsInput struct {
	Query string `json:"query" jsonschema:"description=The query string to search in internal documentation for relevant information and processing steps"`
}

func NewQueryInternalDocsTool() tool.InvokableTool {
	t, err := utils.InferOptionableTool(
		"query_internal_docs",
		`搜索内部知识库和文档。使用向量检索找到相关的操作指南、最佳实践和处理流程。

使用场景:
- 用户问"怎么处理XXX" → query="XXX处理流程"
- 用户问"有什么最佳实践" → query="最佳实践"
- 用户需要操作指导时使用

返回: 相关文档片段列表，按相关性排序`,
		func(ctx context.Context, input *QueryInternalDocsInput, opts ...tool.Option) (output string, err error) {
			rr, err := retriever.NewMilvusRetriever(ctx)
			if err != nil {
				return fmt.Sprintf("retriever error: %v", err), nil
			}
			resp, err := rr.Retrieve(ctx, input.Query)
			if err != nil {
				return fmt.Sprintf("retrieve error: %v", err), nil
			}
			respBytes, _ := json.Marshal(resp)
			return string(respBytes), nil
		})
	if err != nil {
		return nil
	}
	return t
}
