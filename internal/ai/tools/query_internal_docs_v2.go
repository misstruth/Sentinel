package tools

import (
	"SuperBizAgent/internal/ai/retriever"
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

type QueryInternalDocsInputV2 struct {
	Query string `json:"query" jsonschema:"description=The query string to search in internal documentation"`
}

// NewQueryInternalDocsToolV2 创建新版本的知识库检索工具
func NewQueryInternalDocsToolV2(llm model.ChatModel) tool.InvokableTool {
	t, err := utils.InferOptionableTool(
		"query_internal_docs",
		`搜索内部知识库和文档。使用多路召回和重排序找到最相关的文档。

使用场景:
- 用户问"怎么处理XXX" → query="XXX处理流程"
- 用户问"有什么最佳实践" → query="最佳实践"
- 用户需要操作指导时使用

返回: 相关文档片段列表，按相关性排序`,
		func(ctx context.Context, input *QueryInternalDocsInputV2, opts ...tool.Option) (output string, err error) {
			// 使用新的高级检索器
			config := retriever.DefaultConfig()
			config.QueryRewrite.Enabled = true
			config.Rerank.Enabled = false // 如果没有部署 reranker 服务，设为 false

			advRetriever, err := retriever.NewAdvancedRetriever(ctx, llm, config)
			if err != nil {
				return fmt.Sprintf("创建检索器失败: %v", err), nil
			}

			docs, err := advRetriever.Retrieve(ctx, input.Query)
			if err != nil {
				return fmt.Sprintf("检索失败: %v", err), nil
			}

			respBytes, _ := json.Marshal(docs)
			return string(respBytes), nil
		})
	if err != nil {
		return nil
	}
	return t
}
