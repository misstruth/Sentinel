package main

import (
	"SuperBizAgent/internal/ai/evaluation"
	"SuperBizAgent/internal/ai/retriever"
	"context"
	"flag"
	"log"
)

func main() {
	datasetPath := flag.String("dataset", "internal/ai/evaluation/testdata/test_dataset.json", "测试数据集路径")
	outputPath := flag.String("output", "evaluation_report.json", "评测报告输出路径")
	flag.Parse()

	ctx := context.Background()

	// 加载测试数据集
	dataset, err := evaluation.LoadTestDataset(*datasetPath)
	if err != nil {
		log.Fatalf("加载测试数据集失败: %v", err)
	}

	log.Printf("加载测试数据集成功，共 %d 个测试用例", len(dataset.TestCases))

	// 创建检索器
	config := retriever.DefaultConfig()
	config.QueryRewrite.Enabled = true
	config.Rerank.Enabled = false // 暂时禁用 reranker

	r, err := retriever.NewHybridRetriever(ctx, config.HybridSearch.TopK)
	if err != nil {
		log.Fatalf("创建检索器失败: %v", err)
	}

	// 执行评测
	log.Println("开始评测...")
	evaluator := evaluation.NewEvaluator(r)
	metrics, results, err := evaluator.Evaluate(ctx, dataset)
	if err != nil {
		log.Fatalf("评测失败: %v", err)
	}

	// 打印指标
	evaluation.PrintMetrics(metrics)

	// 保存报告
	if err := evaluation.SaveEvaluationReport(metrics, results, *outputPath); err != nil {
		log.Fatalf("保存报告失败: %v", err)
	}

	log.Printf("评测报告已保存到: %s", *outputPath)
}
