#!/bin/bash

# RAG 系统部署脚本

echo "=== RAG 系统部署 ==="

# 1. 检查依赖
echo "检查依赖..."
command -v docker >/dev/null 2>&1 || { echo "需要安装 Docker"; exit 1; }

# 2. 部署 Redis（缓存）
echo "部署 Redis..."
docker run -d \
  --name rag-redis \
  -p 6379:6379 \
  redis:7-alpine

# 3. 部署 BGE Reranker（可选）
echo "部署 BGE Reranker..."
echo "注意: 需要 GPU 支持，如果没有 GPU 可以跳过此步骤"
# docker run -d \
#   --name rag-reranker \
#   -p 8001:8000 \
#   --gpus all \
#   -v /models:/models \
#   reranker-service:latest

# 4. 配置环境变量
echo "配置环境变量..."
export RAG_REDIS_ADDR="localhost:6379"
export RAG_RERANK_URL="http://localhost:8001/rerank"

# 5. 运行评测
echo "运行评测..."
go run internal/ai/cmd/evaluation_cmd/main.go \
  -dataset internal/ai/evaluation/testdata/test_dataset.json \
  -output evaluation_report.json

echo "=== 部署完成 ==="
echo "Redis: localhost:6379"
echo "Reranker: localhost:8001"
echo "评测报告: evaluation_report.json"
