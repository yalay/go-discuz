package ecms

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type EcmsSql struct {
	db *sql.DB
}

// "user:password@/dbname"
func NewEcmsSql(userName, password, database string) *EcmsSql {
	db, err := sql.Open("mysql", userName+":"+password+"@/"+database+"?charset=utf8")
	if err != nil {
		fmt.Printf("new sql err:%v", err)
		return nil
	}
	return &EcmsSql{
		db: db,
	}
}

func (ecms *EcmsSql) GetAllArticle() []*Article {
	rows, err := ecms.db.Query("SELECT a.id, a.classid, a.title, a.ftitle, b.newstext FROM phome_ecms_news a LEFT JOIN phome_ecms_news_data_1 b on a.id=b.id")
	if err != nil {
		fmt.Printf("query err:%v", err)
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
		article.FormatBodyForDiscuz()
		articles = append(articles, article)
	}
	return articles
}

func (ecms *EcmsSql) GetHundredArticles(startId int) ([]*Article, int) {
	querySql := fmt.Sprintf("SELECT a.id, a.classid, a.title, a.ftitle, b.newstext FROM phome_ecms_news a LEFT JOIN phome_ecms_news_data_1 b on a.id=b.id where a.id > %d LIMIT 100", startId)
	rows, err := ecms.db.Query(querySql)
	if err != nil {
		fmt.Printf("query err:%v", err)
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
		article.FormatBodyForDiscuz()
		articles = append(articles, article)
	}
	return articles, maxId
}

func (ecms *EcmsSql) Close() {
	ecms.db.Close()
}
