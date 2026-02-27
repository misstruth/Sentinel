package event_indexer

import (
	"SuperBizAgent/utility/client"
	"SuperBizAgent/utility/common"
	"context"
	"log"
	"strconv"
)

const maxDocuments = 10000

// EnsureCapacity 检查Milvus集合容量并flush确保数据持久化
func EnsureCapacity(ctx context.Context) {
	cli, err := client.NewMilvusClient(ctx)
	if err != nil {
		log.Printf("[event_indexer] milvus client error: %v", err)
		return
	}

	stats, err := cli.GetCollectionStatistics(ctx, common.MilvusCollectionName)
	if err != nil {
		log.Printf("[event_indexer] get stats error: %v", err)
		return
	}

	count := 0
	if v, ok := stats["row_count"]; ok {
		count, _ = strconv.Atoi(v)
	}
	log.Printf("[event_indexer] collection %s has %d rows (limit: %d)", common.MilvusCollectionName, count, maxDocuments)

	// flush确保已写入数据持久化
	if err := cli.Flush(ctx, common.MilvusCollectionName, false); err != nil {
		log.Printf("[event_indexer] flush error: %v", err)
	}
}
