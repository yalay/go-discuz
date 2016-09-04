package main

import (
	"discuz"
	"ecms"
	"flag"
	"tools"
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "c", "config.json", "config file")
	flag.Parse()
}

func main() {
	config := tools.LoadConfig(configFile)
	if config == nil {
		return
	}

	// 从ecms的数据库中读取文章，生成文件
	if config.EnableEcms {
		ecms.GenArticleFile(config.DataFile, config.EcmsSqlConfig)
	}

	discuz.PublishArticleFromFile(config)
}
