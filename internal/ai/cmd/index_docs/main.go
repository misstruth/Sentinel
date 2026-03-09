package main

import (
	"SuperBizAgent/internal/ai/indexer"
	"SuperBizAgent/internal/ai/loader"
	"context"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/components/document"
)

func main() {
	ctx := context.Background()

	// 创建索引器
	fmt.Println("创建索引器...")
	idx, err := indexer.NewMilvusIndexer(ctx)
	if err != nil {
		log.Fatalf("创建索引器失败: %v", err)
	}
	fmt.Println("✅ 索引器创建成功")

	// 创建加载器
	fileLoader, err := loader.NewFileLoader(ctx)
	if err != nil {
		log.Fatalf("创建加载器失败: %v", err)
	}

	// 遍历 docs 目录
	count := 0
	err = filepath.WalkDir("./docs", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		fmt.Printf("索引文件: %s\n", path)
		docs, err := fileLoader.Load(ctx, document.Source{URI: path})
		if err != nil {
			log.Printf("  ❌ 加载失败: %v", err)
			return nil
		}

		_, err = idx.Store(ctx, docs)
		if err != nil {
			log.Printf("  ❌ 索引失败: %v", err)
			return nil
		}

		count += len(docs)
		fmt.Printf("  ✅ 已索引 %d 个文档片段\n", len(docs))
		return nil
	})

	if err != nil {
		log.Fatalf("遍历失败: %v", err)
	}

	fmt.Printf("\n✅ 索引完成，共 %d 个文档片段\n", count)
}
