package retriever

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

// QueryRewriter Query 改写器
type QueryRewriter struct {
	llm        model.ChatModel
	numRewrites int
	useHyDE    bool
}

func NewQueryRewriter(llm model.ChatModel, config QueryRewriteConfig) *QueryRewriter {
	return &QueryRewriter{
		llm:        llm,
		numRewrites: config.NumRewrites,
		useHyDE:    config.UseHyDE,
	}
}

func (q *QueryRewriter) Rewrite(ctx context.Context, query string) ([]string, error) {
	if !q.useHyDE {
		return q.rewriteQueries(ctx, query)
	}
	return q.hydeRewrite(ctx, query)
}

func (q *QueryRewriter) rewriteQueries(ctx context.Context, query string) ([]string, error) {
	prompt := fmt.Sprintf(`将用户查询改写为 %d 个更精确的检索查询。

原始查询：%s

要求：
1. 补充领域术语
2. 扩展同义词
3. 明确意图

直接输出改写后的查询，每行一个，不要编号。`, q.numRewrites, query)

	resp, err := q.llm.Generate(ctx, []*schema.Message{
		schema.UserMessage(prompt),
	})
	if err != nil {
		return []string{query}, nil // 降级返回原查询
	}

	// 解析改写结果
	rewrites := parseRewrites(resp.Content)
	if len(rewrites) == 0 {
		return []string{query}, nil
	}

	return append([]string{query}, rewrites...), nil
}

func (q *QueryRewriter) hydeRewrite(ctx context.Context, query string) ([]string, error) {
	prompt := fmt.Sprintf(`假设你是运维专家，请回答以下问题：

%s

要求：给出简洁的答案（100字以内）。`, query)

	resp, err := q.llm.Generate(ctx, []*schema.Message{
		schema.UserMessage(prompt),
	})
	if err != nil {
		return []string{query}, nil
	}

	// 用假设答案作为检索查询
	return []string{query, resp.Content}, nil
}

func parseRewrites(content string) []string {
	lines := []string{}
	for _, line := range splitLines(content) {
		line = trimPrefix(line)
		if len(line) > 0 {
			lines = append(lines, line)
		}
	}
	return lines
}

func splitLines(s string) []string {
	result := []string{}
	current := ""
	for _, c := range s {
		if c == '\n' {
			if len(current) > 0 {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if len(current) > 0 {
		result = append(result, current)
	}
	return result
}

func trimPrefix(s string) string {
	// 移除编号前缀如 "1. ", "- " 等
	prefixes := []string{"1. ", "2. ", "3. ", "- ", "* "}
	for _, prefix := range prefixes {
		if len(s) > len(prefix) && s[:len(prefix)] == prefix {
			return s[len(prefix):]
		}
	}
	return s
}
