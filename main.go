package main

import (
	"database/sql"
	"flag"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var (
	userName string
	password string
	database string
)

type Article struct {
	classId  int
	title    string
	keywords string
	body     string
}

func init() {
	flag.StringVar(&userName, "u", "", "mysql user name")
	flag.StringVar(&password, "p", "", "mysql password")
	flag.StringVar(&database, "d", "", "mysql database name")
	flag.Parse()
}

// "user:password@/dbname"
func main() {
	db, err := sql.Open("mysql", userName+":"+password+"@/"+database+"?charset=utf8")
	checkErr(err)

	//查询数据
	rows, err := db.Query("SELECT a.classid, a.title, a.ftitle, b.newstext FROM phome_ecms_news a LEFT JOIN phome_ecms_news_data_1 b on a.id=b.id LIMIT 10")
	checkErr(err)

	for rows.Next() {
		var article = Article{}
		err = rows.Scan(&article.classId, &article.title, &article.keywords, &article.body)
		checkErr(err)
		fmt.Printf("%+v\n", article)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
