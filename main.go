package main

import (
	"ecms"
	"flag"
	"fmt"
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

// "user:password@/dbname"
func main() {
	ecmsSql := ecms.NewEcmsSql(userName, password, database)
	if ecmsSql == nil {
		return
	}

	articles, _ := ecmsSql.GetHundredArticles(0)
	for _, article := range articles {
		fmt.Println(article.String())
	}
}
