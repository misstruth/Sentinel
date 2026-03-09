from fastapi import FastAPI
from pydantic import BaseModel
from transformers import AutoModelForSequenceClassification, AutoTokenizer
import torch
from typing import List

app = FastAPI()

# 加载模型
model_name = "BAAI/bge-reranker-large"
tokenizer = AutoTokenizer.from_pretrained(model_name)
model = AutoModelForSequenceClassification.from_pretrained(model_name)
model.eval()

if torch.cuda.is_available():
    model = model.cuda()

class RerankRequest(BaseModel):
    query: str
    documents: List[str]

class RerankResponse(BaseModel):
    scores: List[float]

@app.post("/rerank", response_model=RerankResponse)
async def rerank(request: RerankRequest):
    pairs = [[request.query, doc] for doc in request.documents]

    with torch.no_grad():
        inputs = tokenizer(pairs, padding=True, truncation=True,
                          return_tensors='pt', max_length=512)
        if torch.cuda.is_available():
            inputs = {k: v.cuda() for k, v in inputs.items()}

        scores = model(**inputs, return_dict=True).logits.view(-1).float()
        scores = torch.sigmoid(scores).cpu().tolist()

    return RerankResponse(scores=scores)

@app.get("/health")
async def health():
    return {"status": "ok"}
