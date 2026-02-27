package cmd

import "fmt"

// BatchCmd 批量命令
type BatchCmd struct{}

// Run 执行批量操作
func (c *BatchCmd) Run(args []string) {
	if len(args) == 0 {
		fmt.Println("用法: batch <操作>")
		return
	}
}
