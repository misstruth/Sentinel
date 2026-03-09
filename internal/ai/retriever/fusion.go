package retriever

import (
	"sort"
)

// RRFFusion RRF 融合算法
type RRFFusion struct {
	k int // RRF 参数，通常为 60
}

func NewRRFFusion() *RRFFusion {
	return &RRFFusion{k: 60}
}

// Fuse 融合多路召回结果
func (r *RRFFusion) Fuse(results [][]Document) []Document {
	scores := make(map[string]*Document)
	rrfScores := make(map[string]float64)

	// 计算 RRF 分数
	for _, docList := range results {
		for rank, doc := range docList {
			if _, exists := scores[doc.ID]; !exists {
				docCopy := doc
				scores[doc.ID] = &docCopy
				rrfScores[doc.ID] = 0
			}
			rrfScores[doc.ID] += 1.0 / float64(r.k+rank)
		}
	}

	// 转换为列表并排序
	merged := make([]Document, 0, len(scores))
	for id, doc := range scores {
		doc.Score = rrfScores[id]
		merged = append(merged, *doc)
	}

	sort.Slice(merged, func(i, j int) bool {
		return merged[i].Score > merged[j].Score
	})

	return merged
}

// Deduplicate 去重
func Deduplicate(docs []Document) []Document {
	seen := make(map[string]bool)
	result := []Document{}

	for _, doc := range docs {
		if !seen[doc.ID] {
			seen[doc.ID] = true
			result = append(result, doc)
		}
	}

	return result
}
