package main

import (
	"SuperBizAgent/internal/ai/retriever"
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// 性能测试工具
func main() {
	ctx := context.Background()

	// 创建检索器
	config := retriever.DefaultConfig()
	r, err := retriever.NewHybridRetriever(ctx, config.HybridSearch.TopK)
	if err != nil {
		log.Fatal(err)
	}

	// 测试查询
	queries := []string{
		"服务下线怎么处理？",
		"CPU 使用率过高如何排查？",
		"数据库连接池满了怎么办？",
		"如何分析慢查询？",
		"内存泄漏如何定位？",
	}

	// 并发测试
	concurrency := []int{1, 10, 50, 100}

	for _, c := range concurrency {
		fmt.Printf("\n=== 并发数: %d ===\n", c)
		testConcurrency(ctx, r, queries, c)
	}
}

func testConcurrency(ctx context.Context, r retriever.Retriever, queries []string, concurrency int) {
	var wg sync.WaitGroup
	latencies := make([]time.Duration, 0, concurrency*len(queries))
	var mu sync.Mutex

	start := time.Now()

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for _, query := range queries {
				reqStart := time.Now()
				_, err := r.Retrieve(ctx, query)
				latency := time.Since(reqStart)

				mu.Lock()
				latencies = append(latencies, latency)
				mu.Unlock()

				if err != nil {
					log.Printf("Error: %v", err)
				}
			}
		}(i)
	}

	wg.Wait()
	totalTime := time.Since(start)

	// 计算统计
	printStats(latencies, totalTime, concurrency*len(queries))
}

func printStats(latencies []time.Duration, totalTime time.Duration, totalReqs int) {
	if len(latencies) == 0 {
		return
	}

	// 排序
	for i := 0; i < len(latencies); i++ {
		for j := i + 1; j < len(latencies); j++ {
			if latencies[i] > latencies[j] {
				latencies[i], latencies[j] = latencies[j], latencies[i]
			}
		}
	}

	sum := time.Duration(0)
	for _, l := range latencies {
		sum += l
	}

	avg := sum / time.Duration(len(latencies))
	p50 := latencies[len(latencies)*50/100]
	p95 := latencies[len(latencies)*95/100]
	p99 := latencies[len(latencies)*99/100]

	qps := float64(totalReqs) / totalTime.Seconds()

	fmt.Printf("总请求数: %d\n", totalReqs)
	fmt.Printf("总耗时: %v\n", totalTime)
	fmt.Printf("QPS: %.2f\n", qps)
	fmt.Printf("平均延迟: %v\n", avg)
	fmt.Printf("P50: %v\n", p50)
	fmt.Printf("P95: %v\n", p95)
	fmt.Printf("P99: %v\n", p99)
}
