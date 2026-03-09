FROM python:3.10-slim

WORKDIR /app

# 安装依赖
RUN pip install --no-cache-dir \
    torch \
    transformers \
    fastapi \
    uvicorn \
    pydantic

# 下载模型（可选，也可以挂载）
# RUN python -c "from transformers import AutoModel; AutoModel.from_pretrained('BAAI/bge-reranker-large')"

COPY reranker_service.py .

EXPOSE 8000

CMD ["uvicorn", "reranker_service:app", "--host", "0.0.0.0", "--port", "8000"]
