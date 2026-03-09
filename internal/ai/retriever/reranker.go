package retriever

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"
)

// MultiStageReranker 多阶段重排序器
type MultiStageReranker struct {
	config     RerankConfig
	httpClient *http.Client
}

func NewMultiStageReranker(config RerankConfig) *MultiStageReranker {
	return &MultiStageReranker{
		config: config,
		httpClient: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Millisecond,
		},
	}
}

func (m *MultiStageReranker) Rerank(ctx context.Context, query string, docs []Document) ([]Document, error) {
	if len(docs) == 0 {
		return docs, nil
	}

	// Stage 1: 粗排（简单排序，保留 Top30）
	coarseDocs := m.coarseRank(docs, m.config.CoarseTopK)

	// Stage 2: 精排（BGE Reranker，保留 Top10）
	fineDocs, err := m.fineRank(ctx, query, coarseDocs, m.config.FineTopK)
	if err != nil {
		return coarseDocs[:min(m.config.FinalTopK, len(coarseDocs))], nil // 降级
	}

	// Stage 3: 过滤低分文档
	filtered := m.filterByThreshold(fineDocs, m.config.ScoreThreshold)

	// 返回最终 TopK
	return filtered[:min(m.config.FinalTopK, len(filtered))], nil
}

func (m *MultiStageReranker) coarseRank(docs []Document, topK int) []Document {
	// 按原始分数排序
	sorted := make([]Document, len(docs))
	copy(sorted, docs)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Score > sorted[j].Score
	})
	return sorted[:min(topK, len(sorted))]
}

func (m *MultiStageReranker) fineRank(ctx context.Context, query string, docs []Document, topK int) ([]Document, error) {
	if m.config.RerankURL == "" {
		return docs, nil // 未配置 reranker，跳过
	}

	// 调用 BGE Reranker 服务
	scores, err := m.callRerankerService(ctx, query, docs)
	if err != nil {
		return nil, err
	}

	// 更新分数并排序
	for i := range docs {
		if i < len(scores) {
			docs[i].Score = scores[i]
		}
	}

	sort.Slice(docs, func(i, j int) bool {
		return docs[i].Score > docs[j].Score
	})

	return docs[:min(topK, len(docs))], nil
}

func (m *MultiStageReranker) callRerankerService(ctx context.Context, query string, docs []Document) ([]float64, error) {
	texts := make([]string, len(docs))
	for i, doc := range docs {
		texts[i] = doc.Content
	}

	reqBody := map[string]interface{}{
		"query":     query,
		"documents": texts,
	}

	body, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, "POST", m.config.RerankURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("reranker service error: %d", resp.StatusCode)
	}

	var result struct {
		Scores []float64 `json:"scores"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Scores, nil
}

func (m *MultiStageReranker) filterByThreshold(docs []Document, threshold float64) []Document {
	filtered := []Document{}
	for _, doc := range docs {
		if doc.Score >= threshold {
			filtered = append(filtered, doc)
		}
	}
	return filtered
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
