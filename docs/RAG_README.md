# RAG 系统升级

生产级 RAG 检索系统，支持多路召回、重排序、缓存和评测。

## 特性

- ✅ 混合向量检索（Milvus）
- ✅ Query 改写（LLM 驱动）
- ✅ 多阶段重排序（粗排 → 精排 → 过滤）
- ✅ Redis 缓存
- ✅ 降级策略
- ✅ 评测框架（Recall@K、Precision@K、MRR）
- ✅ Prometheus 监控

## 快速开始

```go
import "SuperBizAgent/internal/ai/retriever"

config := retriever.DefaultConfig()
r, _ := retriever.NewHybridRetriever(ctx, config.HybridSearch.TopK)
docs, _ := r.Retrieve(ctx, "服务下线怎么处理？")
```

## 部署

```bash
# 启动 Redis
docker run -d -p 6379:6379 redis:7-alpine

# 运行评测
go run internal/ai/cmd/evaluation_cmd/main.go
```

## 文档

- [技术方案](docs/rag_upgrade_plan.md)
- [进阶方案](docs/rag_advanced_solutions.md)
- [使用指南](docs/rag_usage_guide.md)
- [实施总结](docs/rag_implementation_summary.md)

## 性能指标

| 指标 | 目标 |
|------|------|
| Recall@10 | 90%+ |
| Precision@3 | 85%+ |
| P95 延迟 | <500ms |
| 缓存命中率 | 40%+ |

## 架构

```
Query → 改写 → 并行召回 → RRF融合 → 重排序 → 过滤 → 结果
         ↓                                        ↑
       缓存检查 ←──────────────────────────── 写入缓存
```
