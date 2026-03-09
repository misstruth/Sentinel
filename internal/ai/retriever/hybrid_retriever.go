package retriever

import (
	"SuperBizAgent/internal/ai/embedder"
	"SuperBizAgent/utility/client"
	"SuperBizAgent/utility/common"
	"context"

	"github.com/cloudwego/eino-ext/components/retriever/milvus"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

// HybridRetriever 混合向量检索器
type HybridRetriever struct {
	milvusRetriever retriever.Retriever
	topK            int
}

func NewHybridRetriever(ctx context.Context, topK int) (*HybridRetriever, error) {
	cli, err := client.NewMilvusClient(ctx)
	if err != nil {
		return nil, err
	}

	eb, err := embedder.DoubaoEmbedding(ctx)
	if err != nil {
		return nil, err
	}

	r, err := milvus.NewRetriever(ctx, &milvus.RetrieverConfig{
		Client:      cli,
		Collection:  common.MilvusCollectionName,
		VectorField: "vector",
		TopK:        topK,
		Embedding:   eb,
		MetricType:  entity.L2,
	})
	if err != nil {
		return nil, err
	}

	return &HybridRetriever{
		milvusRetriever: r,
		topK:            topK,
	}, nil
}

func (h *HybridRetriever) Retrieve(ctx context.Context, query string) ([]Document, error) {
	docs, err := h.milvusRetriever.Retrieve(ctx, query)
	if err != nil {
		return nil, err
	}

	// 转换为内部 Document 格式
	result := make([]Document, len(docs))
	for i, doc := range docs {
		result[i] = Document{
			ID:       doc.ID,
			Content:  doc.Content,
			Score:    doc.Score(),
			Metadata: doc.MetaData,
		}
	}

	return result, nil
}
