# RAG 系统升级使用指南

## 快速开始

### 1. 基础使用

```go
package main

import (
    "SuperBizAgent/internal/ai/retriever"
    "context"
)

func main() {
    ctx := context.Background()

    // 创建检索器
    config := retriever.DefaultConfig()
    r, _ := retriever.NewHybridRetriever(ctx, config.HybridSearch.TopK)

    // 检索
    docs, _ := r.Retrieve(ctx, "服务下线怎么处理？")

    // 使用结果
    for _, doc := range docs {
        println(doc.Content)
    }
}
```

### 2. 启用 Query 改写

```go
config := retriever.DefaultConfig()
config.QueryRewrite.Enabled = true
config.QueryRewrite.NumRewrites = 3

advRetriever, _ := retriever.NewAdvancedRetriever(ctx, llm, config)
docs, _ := advRetriever.Retrieve(ctx, query)
```

### 3. 启用缓存

```go
config.Cache.Enabled = true
config.Cache.TTL = time.Hour

// 需要先启动 Redis
// docker run -d -p 6379:6379 redis:7-alpine
```

### 4. 启用重排序

```go
config.Rerank.Enabled = true
config.Rerank.RerankURL = "http://localhost:8001/rerank"
```

## 部署

### 方式1: 使用脚本部署

```bash
chmod +x scripts/deploy_rag.sh
./scripts/deploy_rag.sh
```

### 方式2: 手动部署

```bash
# 1. 启动 Redis
docker run -d --name rag-redis -p 6379:6379 redis:7-alpine

# 2. 启动 Reranker（可选，需要 GPU）
docker run -d --name rag-reranker -p 8001:8000 --gpus all reranker:latest

# 3. 运行评测
go run internal/ai/cmd/evaluation_cmd/main.go
```

## 评测

### 运行评测

```bash
go run internal/ai/cmd/evaluation_cmd/main.go \
  -dataset internal/ai/evaluation/testdata/test_dataset.json \
  -output evaluation_report.json
```

### 查看结果

```bash
cat evaluation_report.json
```

## 集成到现有 Pipeline

### 替换旧的检索器

```go
// 旧代码
rr, _ := retriever.NewMilvusRetriever(ctx)

// 新代码
config := retriever.DefaultConfig()
rr, _ := retriever.NewAdvancedRetriever(ctx, llm, config)
```

### 更新 Tool

```go
// 使用新版本 Tool
config.ToolsConfig.Tools = append(
    config.ToolsConfig.Tools,
    tools.NewQueryInternalDocsToolV2(llm),
)
```

## 监控

### Prometheus 指标

- `rag_retrieval_latency_ms`: 检索延迟
- `rag_cache_hits_total`: 缓存命中数
- `rag_retrieval_total`: 检索总数

### Grafana 面板

导入 `configs/grafana_dashboard.json`

## 配置说明

### 完整配置

参考 `configs/rag_config.yaml`

### 关键参数

- `hybrid_search.top_k`: 召回数量（默认 20）
- `query_rewrite.num_rewrites`: 改写数量（默认 3）
- `rerank.final_top_k`: 最终返回数量（默认 3）
- `cache.ttl`: 缓存过期时间（默认 1小时）

## 故障排查

### 问题1: 检索延迟高

- 检查 Milvus 性能
- 减少 `top_k` 数量
- 启用缓存

### 问题2: 召回质量差

- 增加 `top_k`
- 启用 Query 改写
- 检查文档索引质量

### 问题3: Reranker 超时

- 增加 `rerank.timeout`
- 检查 Reranker 服务状态
- 临时禁用重排序

## 性能优化建议

1. **启用缓存**: 可提升 40% 命中率
2. **调整 TopK**: 根据场景调整召回数量
3. **批量检索**: 合并多个查询
4. **异步处理**: 使用 goroutine 并行召回
