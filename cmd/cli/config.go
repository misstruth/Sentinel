package cmd

import (
	"flag"
)

// Config CLI 配置
type Config struct {
	ConfigFile string
	Verbose    bool
}

// ParseFlags 解析命令行参数
func ParseFlags() *Config {
	cfg := &Config{}
	flag.StringVar(&cfg.ConfigFile, "c", "", "配置文件")
	flag.BoolVar(&cfg.Verbose, "v", false, "详细输出")
	flag.Parse()
	return cfg
}
