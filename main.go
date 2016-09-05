package main

import (
	"discuz"
	"ecms"
	"flag"
	"fmt"
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

	fmt.Printf("config:%+v\n", config)
	// 从ecms的数据库中读取文章，生成文件
	if config.EnableEcms {
		ecms.GenArticleFile(config.DataFile, config.EcmsSql)
	}

	// 从数据文件恢复到数据库中
	if config.EnableDiscuz {
		discuz.PublishArticleFromFile(config)
	}
}
