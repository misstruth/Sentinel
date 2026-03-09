package evaluation

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// LoadTestDataset 加载测试数据集
func LoadTestDataset(path string) (*TestDataset, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var dataset TestDataset
	if err := json.Unmarshal(data, &dataset); err != nil {
		return nil, err
	}

	return &dataset, nil
}

// SaveEvaluationReport 保存评测报告
func SaveEvaluationReport(metrics *EvaluationMetrics, results []EvaluationResult, outputPath string) error {
	report := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"metrics":   metrics,
		"results":   results,
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, data, 0644)
}

// PrintMetrics 打印评测指标
func PrintMetrics(metrics *EvaluationMetrics) {
	fmt.Println("=== Evaluation Metrics ===")
	fmt.Printf("Recall@3:     %.2f%%\n", metrics.RecallAt3*100)
	fmt.Printf("Recall@10:    %.2f%%\n", metrics.RecallAt10*100)
	fmt.Printf("Precision@3:  %.2f%%\n", metrics.PrecisionAt3*100)
	fmt.Printf("MRR:          %.4f\n", metrics.MRR)
	fmt.Println("==========================")
}
