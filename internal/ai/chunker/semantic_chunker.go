package chunker

import (
	"context"
	"strings"

	"github.com/cloudwego/eino/components/embedding"
)

// SemanticChunker 语义分块器
type SemanticChunker struct {
	embedder  embedding.Embedder
	threshold float64 // 语义相似度阈值
	maxSize   int     // 最大 chunk 大小
	minSize   int     // 最小 chunk 大小
}

type ChunkConfig struct {
	Threshold float64 `json:"threshold"` // 0.7
	MaxSize   int     `json:"max_size"`  // 1024
	MinSize   int     `json:"min_size"`  // 256
}

func NewSemanticChunker(embedder embedding.Embedder, config ChunkConfig) *SemanticChunker {
	if config.Threshold == 0 {
		config.Threshold = 0.7
	}
	if config.MaxSize == 0 {
		config.MaxSize = 1024
	}
	if config.MinSize == 0 {
		config.MinSize = 256
	}

	return &SemanticChunker{
		embedder:  embedder,
		threshold: config.Threshold,
		maxSize:   config.MaxSize,
		minSize:   config.MinSize,
	}
}

// Chunk 语义分块
func (s *SemanticChunker) Chunk(ctx context.Context, text string) ([]string, error) {
	sentences := splitSentences(text)
	if len(sentences) == 0 {
		return []string{}, nil
	}

	chunks := []string{}
	currentChunk := sentences[0]

	for i := 1; i < len(sentences); i++ {
		// 检查当前 chunk 大小
		if len(currentChunk)+len(sentences[i]) > s.maxSize {
			chunks = append(chunks, currentChunk)
			currentChunk = sentences[i]
			continue
		}

		// 计算语义相似度
		sim, err := s.similarity(ctx, currentChunk, sentences[i])
		if err != nil || sim < s.threshold {
			// 语义断裂或计算失败，开始新 chunk
			if len(currentChunk) >= s.minSize {
				chunks = append(chunks, currentChunk)
				currentChunk = sentences[i]
			} else {
				currentChunk += " " + sentences[i]
			}
		} else {
			currentChunk += " " + sentences[i]
		}
	}

	if len(currentChunk) > 0 {
		chunks = append(chunks, currentChunk)
	}

	return chunks, nil
}

func (s *SemanticChunker) similarity(ctx context.Context, text1, text2 string) (float64, error) {
	emb1, err := s.embedder.EmbedStrings(ctx, []string{text1})
	if err != nil {
		return 0, err
	}

	emb2, err := s.embedder.EmbedStrings(ctx, []string{text2})
	if err != nil {
		return 0, err
	}

	return cosineSimilarity(emb1[0], emb2[0]), nil
}

func splitSentences(text string) []string {
	// 简化实现：按句号、问号、感叹号分割
	text = strings.ReplaceAll(text, "。", "。\n")
	text = strings.ReplaceAll(text, "？", "？\n")
	text = strings.ReplaceAll(text, "！", "！\n")
	text = strings.ReplaceAll(text, ". ", ".\n")
	text = strings.ReplaceAll(text, "? ", "?\n")
	text = strings.ReplaceAll(text, "! ", "!\n")

	lines := strings.Split(text, "\n")
	sentences := []string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			sentences = append(sentences, line)
		}
	}
	return sentences
}

func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (sqrt(normA) * sqrt(normB))
}

func sqrt(x float64) float64 {
	if x == 0 {
		return 0
	}
	z := x
	for i := 0; i < 10; i++ {
		z = (z + x/z) / 2
	}
	return z
}
