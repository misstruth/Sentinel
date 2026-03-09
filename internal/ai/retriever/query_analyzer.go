package retriever

import (
	"context"
	"strings"

	"github.com/cloudwego/eino/components/model"
)

// QueryAnalyzerImpl 查询分析器实现
type QueryAnalyzerImpl struct {
	llm model.ChatModel
}

func NewQueryAnalyzer(llm model.ChatModel) *QueryAnalyzerImpl {
	return &QueryAnalyzerImpl{llm: llm}
}

func (q *QueryAnalyzerImpl) Analyze(ctx context.Context, query string) (*QueryAnalysis, error) {
	// 简单规则：根据查询长度判断复杂度
	complexity := "simple"
	topK := 3

	words := strings.Fields(query)
	if len(words) > 10 {
		complexity = "complex"
		topK = 10
	} else if len(words) > 5 {
		complexity = "medium"
		topK = 5
	}

	// 提取关键词（简化版）
	keywords := extractKeywords(query)

	return &QueryAnalysis{
		OriginalQuery: query,
		Keywords:      keywords,
		Complexity:    complexity,
		TopK:          topK,
	}, nil
}

func extractKeywords(query string) []string {
	// 简化实现：分词后过滤停用词
	words := strings.Fields(query)
	stopWords := map[string]bool{
		"的": true, "是": true, "在": true, "了": true,
		"和": true, "有": true, "为": true, "与": true,
	}

	keywords := []string{}
	for _, word := range words {
		if !stopWords[word] && len(word) > 1 {
			keywords = append(keywords, word)
		}
	}
	return keywords
}
