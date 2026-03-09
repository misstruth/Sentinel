#!/bin/bash

echo "=== 部署 BGE Reranker 服务 ==="

cd deployments

# 检查 Docker
if ! command -v docker &> /dev/null; then
    echo "错误: 需要安装 Docker"
    exit 1
fi

# 检查 GPU
if ! command -v nvidia-smi &> /dev/null; then
    echo "警告: 未检测到 NVIDIA GPU，将使用 CPU 模式（速度较慢）"
fi

# 启动服务
echo "启动服务..."
docker-compose up -d

# 等待服务就绪
echo "等待服务启动..."
sleep 10

# 健康检查
echo "健康检查..."
curl -f http://localhost:8001/health || {
    echo "错误: Reranker 服务启动失败"
    docker-compose logs reranker
    exit 1
}

echo "✅ Reranker 服务部署成功"
echo "URL: http://localhost:8001"
