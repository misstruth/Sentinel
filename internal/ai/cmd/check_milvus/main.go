package main

import (
	"SuperBizAgent/utility/client"
	"SuperBizAgent/utility/common"
	"context"
	"fmt"
	"log"
)

func main() {
	ctx := context.Background()

	// 连接 Milvus
	cli, err := client.NewMilvusClient(ctx)
	if err != nil {
		log.Fatalf("连接 Milvus 失败: %v", err)
	}

	// 获取集合统计
	stats, err := cli.GetCollectionStatistics(ctx, common.MilvusCollectionName)
	if err != nil {
		log.Fatalf("获取统计失败: %v", err)
	}

	fmt.Printf("Collection: %s\n", common.MilvusCollectionName)
	fmt.Printf("统计信息: %v\n", stats)

	// 查询前 5 条数据
	expr := ""
	result, err := cli.Query(ctx, common.MilvusCollectionName, []string{}, expr, []string{"id", "content"})
	if err != nil {
		log.Printf("查询失败: %v", err)
	} else {
		fmt.Printf("\n查询返回列数: %d\n", len(result))
		if len(result) > 0 {
			fmt.Printf("第一列名称: %s, 行数: %d\n", result[0].Name(), result[0].Len())
		}
	}
}
