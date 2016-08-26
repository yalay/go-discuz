package main

import (
	"ecms"
	"flag"
	"time"
	"tools"
)

var (
	userName string
	password string
	database string
	log      string
	dict     string
)

func init() {
	flag.StringVar(&userName, "u", "test", "mysql user name")
	flag.StringVar(&password, "p", "test", "mysql password")
	flag.StringVar(&database, "d", "test", "mysql database name")
	flag.StringVar(&dict, "dict", "", "sego dictionary")
	flag.StringVar(&log, "log", "log", "log file")
	flag.Parse()
}

func main() {
	ecmsSql := ecms.NewEcmsSql(userName, password, database)
	if ecmsSql == nil {
		return
	}
	defer ecmsSql.Close()

	keywordsHandler := tools.NewKeywordsHandler(dict)
	startId := 0
	for {
		articles, lastId := ecmsSql.GetHundredArticles(startId)
		if len(articles) == 0 {
			break
		}
		startId = lastId
		for _, article := range articles {
			article.GenKeyWords(keywordsHandler)
			article.Dump(log)
		}
		time.Sleep(time.Second)
	}
}
