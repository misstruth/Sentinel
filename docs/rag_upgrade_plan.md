# Sentinel RAG 系统升级技术方案

## 一、现状分析

### 当前架构
- 单路召回：仅 Milvus 向量检索
- TopK 固定为 3
- 无重排序
- 无评测体系
- 无缓存和降级

### 存在问题
1. 召回覆盖率低，容易漏召回关键文档
2. 排序质量依赖向量相似度，精度不足
3. 无法评估和优化效果
4. 性能和稳定性无保障

---

## 二、升级目标

### 核心指标
- 召回率提升至 90%+（Top10）
- 精排后 Top3 准确率达 85%+
- P95 延迟 < 500ms
- 可用性 99.9%

---

## 三、技术方案

### 3.1 多路召回架构

#### 架构设计
```
用户查询
    ↓
┌─────────────────────────────────┐
│   Query Analyzer (查询分析)      │
│   - 意图识别                     │
│   - 关键词提取                   │
│   - 复杂度评估                   │
└─────────────────────────────────┘
    ↓
┌─────────────────────────────────┐
│   Multi-Recall (并行召回)        │
├─────────────────────────────────┤
│ 1. 向量召回 (Milvus)             │
│    - Dense Embedding             │
│    - TopK = 20                   │
│                                  │
│ 2. 关键词召回 (BM25)             │
│    - Elasticsearch/内存索引      │
│    - TopK = 10                   │
│                                  │
│ 3. 元数据过滤召回                │
│    - 时间范围                    │
│    - 文档类型                    │
│    - 来源系统                    │
└─────────────────────────────────┘
    ↓
┌─────────────────────────────────┐
│   Fusion (结果融合)              │
│   - RRF (Reciprocal Rank Fusion)│
│   - 去重合并                     │
│   - 候选集 ~30 docs              │
└─────────────────────────────────┘
    ↓
┌─────────────────────────────────┐
│   Rerank (重排序)                │
│   - BGE-reranker-large           │
│   - Cross-attention 精排         │
│   - 输出 Top10                   │
└─────────────────────────────────┘
    ↓
┌─────────────────────────────────┐
│   Filter (过滤)                  │
│   - 相似度阈值 > 0.6             │
│   - 动态 TopK 截断               │
│   - 最终返回 Top3-10             │
└─────────────────────────────────┘
```

#### 实现模块

**模块 1：Query Analyzer**
```go
// internal/ai/retriever/query_analyzer.go
type QueryAnalyzer struct {
    complexity ComplexityEstimator  // 查询复杂度评估
    extractor  KeywordExtractor     // 关键词提取
}

type QueryAnalysis struct {
    Keywords   []string
    Complexity string  // "simple" | "medium" | "complex"
    TopK       int     // 动态 TopK
}
```

**模块 2：Multi-Recall Manager**
```go
// internal/ai/retriever/multi_recall.go
type MultiRecallManager struct {
    vectorRecall   VectorRecaller    // Milvus
    keywordRecall  KeywordRecaller   // BM25
    metadataFilter MetadataFilter    // 元数据过滤
}

// 并行召回
func (m *MultiRecallManager) Recall(ctx context.Context, query string) ([]Document, error)
```

**模块 3：Reranker**
```go
// internal/ai/retriever/reranker.go
type Reranker interface {
    Rerank(ctx context.Context, query string, docs []Document) ([]Document, error)
}

type BGEReranker struct {
    client *http.Client
    apiURL string
}
```

---

### 3.2 BM25 关键词召回

#### 技术选型
**方案 A：Elasticsearch（推荐）**
- 优势：成熟稳定，支持中文分词，运维成熟
- 劣势：需要额外部署

**方案 B：内存 BM25（轻量级）**
- 优势：无额外依赖，启动快
- 劣势：数据量大时内存占用高
- 适用：文档量 < 10万

#### 实现方案（Elasticsearch）
```go
// internal/ai/retriever/bm25_retriever.go
type BM25Retriever struct {
    esClient *elasticsearch.Client
    index    string
}

func (r *BM25Retriever) Search(ctx context.Context, query string, topK int) ([]Document, error) {
    // 使用 match_phrase 和 multi_match
    // 中文分词：ik_max_word
}
```

#### 索引结构
```json
{
  "mappings": {
    "properties": {
      "id": {"type": "keyword"},
      "content": {
        "type": "text",
        "analyzer": "ik_max_word",
        "search_analyzer": "ik_smart"
      },
      "metadata": {"type": "object"},
      "timestamp": {"type": "date"}
    }
  }
}
```

---

### 3.3 结果融合策略

#### RRF (Reciprocal Rank Fusion)
```go
// internal/ai/retriever/fusion.go
func RRFFusion(results [][]Document, k int) []Document {
    // RRF Score = Σ 1/(k + rank_i)
    // k = 60 (常用值)
    scores := make(map[string]float64)
    for _, docList := range results {
        for rank, doc := range docList {
            scores[doc.ID] += 1.0 / (60.0 + float64(rank))
        }
    }
    // 按分数排序返回
}
```

---

### 3.4 重排序实现

#### BGE Reranker 部署

**方案 A：本地部署（推荐）**
```bash
# 使用 BAAI/bge-reranker-large
docker run -d \
  -p 8001:8000 \
  -v /models:/models \
  --gpus all \
  reranker-service:latest
```

**方案 B：API 调用**
- 使用豆包或其他云服务的 Reranker API

#### 接口设计
```go
// internal/ai/retriever/reranker.go
type RerankRequest struct {
    Query     string   `json:"query"`
    Documents []string `json:"documents"`
}

type RerankResponse struct {
    Scores []float64 `json:"scores"`
}

func (r *BGEReranker) Rerank(ctx context.Context, query string, docs []Document) ([]Document, error) {
    // 调用 reranker 服务
    // 按 score 重新排序
    // 返回 Top10
}
```

---

### 3.5 数据质量优化

#### 文档切片策略
```go
// internal/ai/indexer/chunker.go
type ChunkConfig struct {
    ChunkSize    int  // 512 tokens
    ChunkOverlap int  // 100 tokens
    Separator    string  // "\n\n"
}

func ChunkDocument(doc string, config ChunkConfig) []string {
    // 按段落切分
    // 保持语义完整性
    // 添加上下文 overlap
}
```

#### 元数据增强
```go
type DocumentMetadata struct {
    Source      string    `json:"source"`       // 来源系统
    DocType     string    `json:"doc_type"`     // 文档类型
    Timestamp   time.Time `json:"timestamp"`    // 时间戳
    Category    string    `json:"category"`     // 分类
    Tags        []string  `json:"tags"`         // 标签
    Importance  int       `json:"importance"`   // 重要度 1-5
}
```

---

### 3.6 评测体系

#### 评测框架设计
```
┌─────────────────────────────────┐
│   Test Dataset (测试集)          │
│   - 100+ 标注问答对              │
│   - Ground Truth 文档 ID         │
└─────────────────────────────────┘
    ↓
┌─────────────────────────────────┐
│   Offline Evaluation (离线评测)  │
├─────────────────────────────────┤
│ 1. Recall@K (召回率)             │
│    - Top3/5/10 命中率            │
│                                  │
│ 2. Precision@K (精确率)          │
│    - 相关文档占比                │
│                                  │
│ 3. MRR (平均倒数排名)            │
│    - 首个相关文档位置            │
│                                  │
│ 4. NDCG (归一化折损累计增益)     │
│    - 考虑排序质量                │
└─────────────────────────────────┘
    ↓
┌─────────────────────────────────┐
│   Hallucination Detection        │
│   - 答案与文档一致性检测         │
│   - NLI 模型判断蕴含关系         │
└─────────────────────────────────┘
```

#### 实现代码
```go
// internal/ai/evaluation/metrics.go
type EvaluationMetrics struct {
    RecallAt3  float64
    RecallAt10 float64
    PrecisionAt3 float64
    MRR        float64
    NDCG       float64
}

func EvaluateRetrieval(testSet []TestCase, retriever Retriever) EvaluationMetrics {
    // 遍历测试集
    // 计算各项指标
}

// internal/ai/evaluation/hallucination.go
func DetectHallucination(answer string, sources []Document) (score float64, err error) {
    // 使用 NLI 模型检测
    // 返回一致性分数 0-1
}
```

#### 测试数据集格式
```json
{
  "test_cases": [
    {
      "id": "001",
      "query": "服务下线怎么处理？",
      "ground_truth_doc_ids": ["doc_123", "doc_456"],
      "expected_answer": "..."
    }
  ]
}
```

---

### 3.7 工程化优化

#### 缓存机制
```go
// internal/ai/retriever/cache.go
type RetrievalCache struct {
    cache *redis.Client
    ttl   time.Duration  // 1小时
}

func (c *RetrievalCache) Get(ctx context.Context, query string) ([]Document, bool) {
    // 查询缓存
}

func (c *RetrievalCache) Set(ctx context.Context, query string, docs []Document) error {
    // 写入缓存
}
```

#### 降级策略
```go
// internal/ai/retriever/fallback.go
type FallbackRetriever struct {
    primary   Retriever  // 主检索器
    secondary Retriever  // 备用检索器（简化版）
}

func (f *FallbackRetriever) Retrieve(ctx context.Context, query string) ([]Document, error) {
    docs, err := f.primary.Retrieve(ctx, query)
    if err != nil {
        // 降级到备用方案
        return f.secondary.Retrieve(ctx, query)
    }
    return docs, nil
}
```

#### 可观测性
```go
// internal/ai/retriever/metrics.go
type RetrievalMetrics struct {
    VectorRecallLatency   prometheus.Histogram
    BM25RecallLatency     prometheus.Histogram
    RerankLatency         prometheus.Histogram
    CacheHitRate          prometheus.Counter
    RecallCount           prometheus.Counter
}

func (m *RetrievalMetrics) RecordRecall(ctx context.Context, stage string, duration time.Duration) {
    // 记录各阶段耗时
}
```

---

## 四、实施计划

### Phase 1：基础设施（Week 1-2）
- [ ] 部署 Elasticsearch 集群
- [ ] 部署 BGE Reranker 服务
- [ ] 配置 Redis 缓存
- [ ] 搭建监控面板

### Phase 2：多路召回（Week 3-4）
- [ ] 实现 BM25 召回模块
- [ ] 实现 RRF 融合逻辑
- [ ] 元数据过滤增强
- [ ] 动态 TopK 调整

### Phase 3：重排序（Week 5）
- [ ] 集成 BGE Reranker
- [ ] 相似度阈值过滤
- [ ] 性能优化

### Phase 4：评测体系（Week 6-7）
- [ ] 构建测试数据集（100+ 样本）
- [ ] 实现评测框架
- [ ] 幻觉检测模块
- [ ] 自动化评测流程

### Phase 5：工程化（Week 8）
- [ ] 缓存机制上线
- [ ] 降级策略实现
- [ ] 监控告警配置
- [ ] 压测和调优

---

## 五、技术栈

### 新增依赖
```go
// go.mod
require (
    github.com/elastic/go-elasticsearch/v8 v8.x.x  // ES 客户端
    github.com/go-redis/redis/v8 v8.x.x            // Redis 缓存
    github.com/prometheus/client_golang v1.x.x     // 监控
)
```

### 外部服务
- Elasticsearch 7.x+（BM25 检索）
- Redis 6.x+（缓存）
- BGE Reranker Service（重排序）
- Prometheus + Grafana（监控）

---

## 六、风险与应对

### 风险 1：Reranker 延迟高
- **应对**：批量请求，异步处理，设置超时降级

### 风险 2：ES 集群故障
- **应对**：降级到纯向量检索，保证基本可用

### 风险 3：缓存穿透
- **应对**：布隆过滤器 + 空值缓存

### 风险 4：评测数据不足
- **应对**：先用 50 个高频问题启动，逐步扩充

---

## 七、预期效果

### 性能指标
| 指标 | 当前 | 目标 | 提升 |
|------|------|------|------|
| Recall@10 | ~60% | 90%+ | +50% |
| Precision@3 | ~70% | 85%+ | +21% |
| P95 延迟 | ~200ms | <500ms | 可控 |
| 缓存命中率 | 0% | 40%+ | 新增 |

### 业务价值
- 减少 Agent 幻觉率 30%+
- 提升用户满意度
- 降低人工介入成本

---

## 八、后续优化方向

1. **Query 改写**：同义词扩展、拼写纠错
2. **混合向量**：稀疏 + 密集向量联合检索
3. **个性化召回**：基于用户历史的偏好调整
4. **主动学习**：根据用户反馈持续优化
