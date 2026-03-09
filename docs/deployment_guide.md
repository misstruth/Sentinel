# RAG 系统完整部署指南

## 一、环境准备

### 必需组件
- Go 1.21+
- Docker & Docker Compose
- Redis 7+
- Milvus（已有）

### 可选组件
- NVIDIA GPU（用于 Reranker）
- Prometheus + Grafana（监控）

---

## 二、快速部署

### 方式1: 一键部署（推荐）

```bash
# 部署所有服务
./scripts/deploy_rag.sh
```

### 方式2: 分步部署

```bash
# 1. 启动 Redis
docker run -d --name rag-redis -p 6379:6379 redis:7-alpine

# 2. 启动 Reranker（可选）
cd deployments
docker-compose up -d reranker

# 3. 验证服务
curl http://localhost:8001/health
```

---

## 三、配置

### 基础配置

编辑 `configs/rag_config.yaml`:

```yaml
query_rewrite:
  enabled: true
  num_rewrites: 3

rerank:
  enabled: false  # 如无 GPU 设为 false
  rerank_url: "http://localhost:8001/rerank"

cache:
  enabled: true
  redis_addr: "localhost:6379"
```

---

## 四、运行测试

### 评测测试

```bash
go run internal/ai/cmd/evaluation_cmd/main.go
```

### 性能测试

```bash
go run internal/ai/cmd/benchmark_cmd/main.go
```

### 功能测试

```bash
go run internal/ai/cmd/rag_example/main.go
```

---

## 五、集成到项目

### 替换旧检索器

```go
// 在 chat_pipeline/flow.go 中
import "SuperBizAgent/internal/ai/tools"

// 替换 Tool
config.ToolsConfig.Tools = append(
    config.ToolsConfig.Tools,
    tools.NewQueryInternalDocsToolV2(chatModelIns11),
)
```

---

## 六、监控

### Prometheus 指标

访问: http://localhost:9090

关键指标:
- `rag_retrieval_latency_ms`
- `rag_cache_hits_total`
- `rag_retrieval_total`

---

## 七、故障排查

### Redis 连接失败
```bash
docker ps | grep redis
docker logs rag-redis
```

### Reranker 超时
```bash
# 检查服务状态
curl http://localhost:8001/health

# 查看日志
docker logs rag-reranker
```

### 检索质量差
```bash
# 运行评测
go run internal/ai/cmd/evaluation_cmd/main.go

# 查看报告
cat evaluation_report.json
```
