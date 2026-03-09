package chunker

import "strings"

// FixedChunker 固定长度分块器
type FixedChunker struct {
	chunkSize int
	overlap   int
}

func NewFixedChunker(chunkSize, overlap int) *FixedChunker {
	return &FixedChunker{
		chunkSize: chunkSize,
		overlap:   overlap,
	}
}

func (f *FixedChunker) Chunk(text string) []string {
	runes := []rune(text)
	chunks := []string{}

	for i := 0; i < len(runes); i += f.chunkSize - f.overlap {
		end := i + f.chunkSize
		if end > len(runes) {
			end = len(runes)
		}

		chunk := string(runes[i:end])
		chunk = strings.TrimSpace(chunk)
		if len(chunk) > 0 {
			chunks = append(chunks, chunk)
		}

		if end >= len(runes) {
			break
		}
	}

	return chunks
}
