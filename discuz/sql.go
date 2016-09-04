package discuz

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var author string
var authorId int
var dbPrefix = "pre_"

type DiscuzSql struct {
	db *sql.DB
}

// "user:password@/dbname"
func newDiscuzSql(userName, password, database string) *DiscuzSql {
	db, err := sql.Open("mysql", userName+":"+password+"@/"+database+"?charset=utf8")
	if err != nil {
		fmt.Printf("new sql err:%v", err)
		return nil
	}
	return &DiscuzSql{
		db: db,
	}
}

// 获取文章页id
func (d *DiscuzSql) GetPostId() int64 {
	sqlPre, err := d.db.Prepare(`INSERT ` + dbPrefix + `forum_post_tableid (pid) VALUES (?)`)
	if err != nil {
		fmt.Printf("Prepare sql err:%v", err)
		return 0
	}
	sqlResp, err := sqlPre.Exec(0)
	if err != nil {
		fmt.Printf("Exec sql err:%v", err)
		return 0
	}

	postId, err := sqlResp.LastInsertId()
	if err != nil {
		fmt.Printf("Get Last insert id err:%v", err)
		return 0
	}

	return postId
}

// 插入thread列表
func (d *DiscuzSql) InsertThread(article *Article) int64 {
	sqlPre, err := d.db.Prepare(`INSERT ` + dbPrefix + `forum_thread (fid, subject, author, authorid, 
		dateline, lastpost, lastposter) VALUES (?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		fmt.Printf("Prepare sql err:%v", err)
		return 0
	}

	dateLine := time.Now().Unix()
	sqlResp, err := sqlPre.Exec(article.ClassId, article.Title, author, authorId, dateLine, dateLine, author)
	if err != nil {
		fmt.Printf("Exec sql err:%v", err)
		return 0
	}

	threadId, err := sqlResp.LastInsertId()
	if err != nil {
		fmt.Printf("Get Last insert id err:%v", err)
		return 0
	}

	return threadId
}

func (d *DiscuzSql) InsertPost(article *Article, pid, tid int64) {
	sqlPre, err := d.db.Prepare(`INSERT ` + dbPrefix + `forum_post (pid, fid, tid, first,
		author, authorId, subject, dateline, message, usesig, smileyoff, position,
		tags) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		fmt.Printf("Prepare sql err:%v", err)
		return
	}

	dateLine := time.Now().Unix()
	sqlResp, err := sqlPre.Exec(pid, article.ClassId, tid, 1, author, authorId, article.Title,
		dateLine, article.Body, 1, -1, 1, article.Keywords)
	if err != nil {
		fmt.Printf("Exec sql err:%v", err)
		return
	}
	fmt.Printf("Exec sql success:%v", sqlResp)
}

func (d *DiscuzSql) UpdateForum(article *Article, tid int64) {
	sqlPre, err := d.db.Prepare(`UPDATE ` + dbPrefix + `forum_forum SET threads=threads+1,
		posts=posts+1, todayposts=todayposts+1, lastpost=? WHERE where fid=?`)
	if err != nil {
		fmt.Printf("Prepare sql err:%v", err)
		return
	}

	lastpost := fmt.Sprintf("%d", tid) + "\t" + article.Title + "\t" + fmt.Sprintf("%d", time.Now().Unix()) + "\t" + author
	sqlResp, err := sqlPre.Exec(lastpost, article.ClassId)
	if err != nil {
		fmt.Printf("Exec sql err:%v", err)
		return
	}
	fmt.Printf("Exec sql success:%v", sqlResp)
}

func (d *DiscuzSql) Close() {
	d.db.Close()
}
