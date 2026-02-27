package cmd

import (
	"fmt"
	"os"
)

// Execute 执行 CLI
func Execute() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "help":
		printHelp()
	case "version":
		printVersion()
	default:
		fmt.Println("未知命令:", os.Args[1])
	}
}

func printHelp() {
	fmt.Println("Fo-Sentinel CLI")
	fmt.Println("用法: fo-sentinel <命令>")
	fmt.Println("")
	fmt.Println("命令:")
	fmt.Println("  help     显示帮助")
	fmt.Println("  version  显示版本")
}

func printVersion() {
	fmt.Println("Fo-Sentinel v1.0.0")
}
