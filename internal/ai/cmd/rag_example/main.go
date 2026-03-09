package main

import (
	"SuperBizAgent/internal/ai/retriever"
	"context"
	"fmt"
	"log"
)

// 使用示例
func main() {
	ctx := context.Background()

	// 方式1: 使用默认配置
	config := retriever.DefaultConfig()

	// 方式2: 自定义配置
	config.QueryRewrite.Enabled = true
	config.QueryRewrite.NumRewrites = 3
	config.Rerank.Enabled = false // 如果没有部署 reranker

	// 创建检索器
	r, err := retriever.NewHybridRetriever(ctx, config.HybridSearch.TopK)
	if err != nil {
		log.Fatal(err)
	}

	// 执行检索
	query := "服务下线怎么处理？"
	docs, err := r.Retrieve(ctx, query)
	if err != nil {
		log.Fatal(err)
	}

	// 打印结果
	fmt.Printf("查询: %s\n", query)
	fmt.Printf("找到 %d 个相关文档:\n", len(docs))
	for i, doc := range docs {
		fmt.Printf("%d. [%.4f] %s\n", i+1, doc.Score, doc.Content[:100])
	}
}
