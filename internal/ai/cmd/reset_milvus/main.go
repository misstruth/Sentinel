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

	cli, err := client.NewMilvusClient(ctx)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}

	// 删除旧集合
	fmt.Printf("删除集合: %s\n", common.MilvusCollectionName)
	err = cli.DropCollection(ctx, common.MilvusCollectionName)
	if err != nil {
		log.Printf("删除集合失败（可能不存在）: %v", err)
	} else {
		fmt.Println("✅ 集合已删除")
	}

	fmt.Println("\n请重新运行程序以创建新集合")
}
