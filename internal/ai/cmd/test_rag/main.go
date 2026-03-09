package main

import (
	"SuperBizAgent/internal/ai/retriever"
	"context"
	"fmt"
	"log"
)

func main() {
	ctx := context.Background()

	// 测试基础检索器
	fmt.Println("=== 测试 1: 创建检索器 ===")
	r, err := retriever.NewHybridRetriever(ctx, 3)
	if err != nil {
		log.Fatalf("创建检索器失败: %v", err)
	}
	fmt.Println("✅ 检索器创建成功")

	// 测试检索
	fmt.Println("\n=== 测试 2: 执行检索 ===")
	query := "服务下线怎么处理？"
	docs, err := r.Retrieve(ctx, query)
	if err != nil {
		log.Printf("❌ 检索失败: %v", err)
	} else {
		fmt.Printf("✅ 检索成功，返回 %d 个文档\n", len(docs))
		for i, doc := range docs {
			fmt.Printf("  %d. [%.4f] %s\n", i+1, doc.Score, truncate(doc.Content, 50))
		}
	}

	// 测试配置
	fmt.Println("\n=== 测试 3: 配置加载 ===")
	config := retriever.DefaultConfig()
	fmt.Printf("✅ 默认配置:\n")
	fmt.Printf("  - TopK: %d\n", config.HybridSearch.TopK)
	fmt.Printf("  - Query改写: %v\n", config.QueryRewrite.Enabled)
	fmt.Printf("  - 重排序: %v\n", config.Rerank.Enabled)
	fmt.Printf("  - 缓存: %v\n", config.Cache.Enabled)

	fmt.Println("\n=== 所有测试完成 ===")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
