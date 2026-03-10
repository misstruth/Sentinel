package fetcher

import (
	"crypto/sha256"
	"fmt"
)

// generateHash 生成唯一哈希
func generateHash(source, id string) string {
	h := sha256.New()
	h.Write([]byte(source + ":" + id))
	return fmt.Sprintf("%x", h.Sum(nil))
}
