package evaluation

// TestCase 测试用例
type TestCase struct {
	ID              string   `json:"id"`
	Query           string   `json:"query"`
	GroundTruthDocs []string `json:"ground_truth_doc_ids"`
	ExpectedAnswer  string   `json:"expected_answer,omitempty"`
	Category        string   `json:"category"`
	Difficulty      string   `json:"difficulty"` // "easy" | "medium" | "hard"
}

// TestDataset 测试数据集
type TestDataset struct {
	Version   string     `json:"version"`
	CreatedAt string     `json:"created_at"`
	TestCases []TestCase `json:"test_cases"`
}

// EvaluationMetrics 评测指标
type EvaluationMetrics struct {
	RecallAt3    float64 `json:"recall_at_3"`
	RecallAt5    float64 `json:"recall_at_5"`
	RecallAt10   float64 `json:"recall_at_10"`
	PrecisionAt3 float64 `json:"precision_at_3"`
	MRR          float64 `json:"mrr"`
	NDCG         float64 `json:"ndcg"`
}

// EvaluationResult 评测结果
type EvaluationResult struct {
	TestCaseID   string   `json:"test_case_id"`
	Query        string   `json:"query"`
	RetrievedDocs []string `json:"retrieved_doc_ids"`
	GroundTruth  []string `json:"ground_truth_doc_ids"`
	RecallAt3    float64  `json:"recall_at_3"`
	RecallAt10   float64  `json:"recall_at_10"`
	MRR          float64  `json:"mrr"`
}
