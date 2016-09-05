package ecms

import (
	"database/sql"
	"fmt"
	"time"
	"tools"

	_ "github.com/go-sql-driver/mysql"
)

type EcmsSql struct {
	db *sql.DB
}

func GenArticleFile(dataFile string, config *tools.SqlConfig) {
	if config == nil {
		return
	}
	ecmsSql := newEcmsSql(config.UserName, config.Password, config.Database)
	if ecmsSql == nil {
		return
	}
	defer ecmsSql.close()

	startId := 0
	for {
		articles, lastId := ecmsSql.getHundredArticles(startId)
		if len(articles) == 0 {
			break
		}
		startId = lastId
		for _, article := range articles {
			article.Dump(dataFile)
		}
		time.Sleep(time.Second)
	}
}

// "user:password@/dbname"
func newEcmsSql(userName, password, database string) *EcmsSql {
	db, err := sql.Open("mysql", userName+":"+password+"@/"+database+"?charset=utf8")
	if err != nil {
		fmt.Printf("new sql err:%v\n", err)
		return nil
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("db connect err:%v\n", err)
		return nil
	}

	return &EcmsSql{
		db: db,
	}
}

func (e *EcmsSql) getAllArticle() []*Article {
	rows, err := e.db.Query("SELECT a.id, a.classid, a.title, a.ftitle, b.newstext FROM phome_ecms_news a LEFT JOIN phome_ecms_news_data_1 b on a.id=b.id")
	if err != nil {
		fmt.Printf("query err:%v\n", err)
		return nil
	}

	articles := make([]*Article, 0)
	for rows.Next() {
		var curId, classId int
		var title, keywords, body string
		rows.Scan(&curId, &classId, &title, &keywords, &body)
		article := NewArticle(curId, classId, title, keywords, body)
		if article == nil {
			continue
		}
		articles = append(articles, article)
	}
	return articles
}

func (e *EcmsSql) getHundredArticles(startId int) ([]*Article, int) {
	querySql := fmt.Sprintf("SELECT a.id, a.classid, a.title, a.ftitle, b.newstext FROM phome_ecms_news a LEFT JOIN phome_ecms_news_data_1 b on a.id=b.id where a.id > %d LIMIT 100", startId)
	rows, err := e.db.Query(querySql)
	if err != nil {
		fmt.Printf("query err:%v\n", err)
		return nil, startId
	}

	maxId := startId
	articles := make([]*Article, 0, 100)
	for rows.Next() {
		var curId, classId int
		var title, keywords, body string
		rows.Scan(&curId, &classId, &title, &keywords, &body)
		if curId > maxId {
			maxId = curId
		}
		article := NewArticle(curId, classId, title, keywords, body)
		if article == nil {
			continue
		}
		articles = append(articles, article)
	}
	return articles, maxId
}

func (e *EcmsSql) close() {
	e.db.Close()
}
