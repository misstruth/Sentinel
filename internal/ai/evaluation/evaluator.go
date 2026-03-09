package evaluation

import (
	"SuperBizAgent/internal/ai/retriever"
	"context"
	"math"
)

// Evaluator 评测器
type Evaluator struct {
	retriever retriever.Retriever
}

func NewEvaluator(r retriever.Retriever) *Evaluator {
	return &Evaluator{retriever: r}
}

// Evaluate 执行评测
func (e *Evaluator) Evaluate(ctx context.Context, dataset *TestDataset) (*EvaluationMetrics, []EvaluationResult, error) {
	results := []EvaluationResult{}
	totalRecall3 := 0.0
	totalRecall10 := 0.0
	totalPrecision3 := 0.0
	totalMRR := 0.0

	for _, testCase := range dataset.TestCases {
		docs, err := e.retriever.Retrieve(ctx, testCase.Query)
		if err != nil {
			continue
		}

		retrievedIDs := make([]string, len(docs))
		for i, doc := range docs {
			retrievedIDs[i] = doc.ID
		}

		// 计算指标
		recall3 := calculateRecall(retrievedIDs[:min(3, len(retrievedIDs))], testCase.GroundTruthDocs)
		recall10 := calculateRecall(retrievedIDs[:min(10, len(retrievedIDs))], testCase.GroundTruthDocs)
		mrr := calculateMRR(retrievedIDs, testCase.GroundTruthDocs)

		results = append(results, EvaluationResult{
			TestCaseID:    testCase.ID,
			Query:         testCase.Query,
			RetrievedDocs: retrievedIDs,
			GroundTruth:   testCase.GroundTruthDocs,
			RecallAt3:     recall3,
			RecallAt10:    recall10,
			MRR:           mrr,
		})

		totalRecall3 += recall3
		totalRecall10 += recall10
		totalMRR += mrr

		if len(retrievedIDs) >= 3 {
			precision3 := calculatePrecision(retrievedIDs[:3], testCase.GroundTruthDocs)
			totalPrecision3 += precision3
		}
	}

	n := float64(len(dataset.TestCases))
	metrics := &EvaluationMetrics{
		RecallAt3:    totalRecall3 / n,
		RecallAt10:   totalRecall10 / n,
		PrecisionAt3: totalPrecision3 / n,
		MRR:          totalMRR / n,
	}

	return metrics, results, nil
}

func calculateRecall(retrieved, groundTruth []string) float64 {
	if len(groundTruth) == 0 {
		return 0
	}

	gtSet := make(map[string]bool)
	for _, id := range groundTruth {
		gtSet[id] = true
	}

	hits := 0
	for _, id := range retrieved {
		if gtSet[id] {
			hits++
		}
	}

	return float64(hits) / float64(len(groundTruth))
}

func calculatePrecision(retrieved, groundTruth []string) float64 {
	if len(retrieved) == 0 {
		return 0
	}

	gtSet := make(map[string]bool)
	for _, id := range groundTruth {
		gtSet[id] = true
	}

	hits := 0
	for _, id := range retrieved {
		if gtSet[id] {
			hits++
		}
	}

	return float64(hits) / float64(len(retrieved))
}

func calculateMRR(retrieved, groundTruth []string) float64 {
	gtSet := make(map[string]bool)
	for _, id := range groundTruth {
		gtSet[id] = true
	}

	for i, id := range retrieved {
		if gtSet[id] {
			return 1.0 / float64(i+1)
		}
	}

	return 0
}

func calculateNDCG(retrieved, groundTruth []string, k int) float64 {
	dcg := 0.0
	for i := 0; i < min(k, len(retrieved)); i++ {
		if contains(groundTruth, retrieved[i]) {
			dcg += 1.0 / math.Log2(float64(i+2))
		}
	}

	idcg := 0.0
	for i := 0; i < min(k, len(groundTruth)); i++ {
		idcg += 1.0 / math.Log2(float64(i+2))
	}

	if idcg == 0 {
		return 0
	}
	return dcg / idcg
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
