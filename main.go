package main

import (
	"ecms"
	"flag"
	"time"
)

var (
	userName string
	password string
	database string
	log      string
)

func init() {
	flag.StringVar(&userName, "u", "test", "mysql user name")
	flag.StringVar(&password, "p", "test", "mysql password")
	flag.StringVar(&database, "d", "test", "mysql database name")
	flag.StringVar(&log, "log", "log", "log file")
	flag.Parse()
}

func main() {
	ecmsSql := ecms.NewEcmsSql(userName, password, database)
	if ecmsSql == nil {
		return
	}
	defer ecmsSql.Close()

	startId := 0
	for {
		articles, lastId := ecmsSql.GetHundredArticles(startId)
		if len(articles) == 0 {
			break
		}
		startId = lastId
		for _, article := range articles {
			article.Dump(log)
		}
		time.Sleep(time.Second)
	}
}
