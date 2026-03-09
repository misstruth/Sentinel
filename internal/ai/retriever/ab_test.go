package retriever

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"log"
)

// ABTestRetriever A/B 测试检索器
type ABTestRetriever struct {
	strategyA Retriever
	strategyB Retriever
	ratio     float64 // B 策略流量占比 0-1
}

func NewABTestRetriever(strategyA, strategyB Retriever, ratio float64) *ABTestRetriever {
	return &ABTestRetriever{
		strategyA: strategyA,
		strategyB: strategyB,
		ratio:     ratio,
	}
}

func (a *ABTestRetriever) Retrieve(ctx context.Context, query string) ([]Document, error) {
	// 根据查询哈希决定使用哪个策略
	if a.shouldUseStrategyB(query) {
		log.Printf("[AB Test] Using Strategy B for query: %s", query)
		docs, err := a.strategyB.Retrieve(ctx, query)
		RecordRetrieval(err == nil)
		return docs, err
	}

	log.Printf("[AB Test] Using Strategy A for query: %s", query)
	docs, err := a.strategyA.Retrieve(ctx, query)
	RecordRetrieval(err == nil)
	return docs, err
}

func (a *ABTestRetriever) shouldUseStrategyB(query string) bool {
	hash := md5.Sum([]byte(query))
	hashStr := hex.EncodeToString(hash[:])

	// 取哈希的前8位转为数字
	var hashNum uint32
	for i := 0; i < 8 && i < len(hashStr); i++ {
		hashNum = hashNum*16 + uint32(hashStr[i])
	}

	// 计算百分比
	percentage := float64(hashNum%100) / 100.0
	return percentage < a.ratio
}
