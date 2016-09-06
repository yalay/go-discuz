package discuz

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
	"tools"

	_ "github.com/go-sql-driver/mysql"
)

type DiscuzSql struct {
	db       *sql.DB
	dbPrefix string
	author   string
	authorId int
}

// "user:password@/dbname"
func newDiscuzSql(sqlCfg *tools.SqlConfig) *DiscuzSql {
	db, err := sql.Open("mysql", sqlCfg.UserName+":"+sqlCfg.Password+"@/"+sqlCfg.Database+"?charset=utf8")
	if err != nil {
		fmt.Printf("new sql err:%v\n", err)
		return nil
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("db connect err:%v\n", err)
		return nil
	}

	return &DiscuzSql{
		db:       db,
		dbPrefix: sqlCfg.DbPrefix,
		author:   sqlCfg.Author,
		authorId: sqlCfg.AuthorId,
	}
}

func (d *DiscuzSql) CheckTitleExist(article *Article) bool {
	querySql := fmt.Sprintf(`SELECT tid FROM `+d.dbPrefix+`forum_thread WHERE subject="%s" LIMIT 1`, article.Title)
	rows, err := d.db.Query(querySql)
	if err != nil {
		fmt.Printf("CheckTitleExist query err:%v\n", err)
		return true
	}
	defer rows.Close()

	for rows.Next() {
		var threadId int64
		rows.Scan(&threadId)
		fmt.Printf("CheckTitleExist err. Title:%s, exist tid:%d\n", article.Title, threadId)
		return true
	}
	return false
}

// 获取文章页id
func (d *DiscuzSql) GetPostId() int64 {
	sqlPre, err := d.db.Prepare(`INSERT ` + d.dbPrefix + `forum_post_tableid (pid) VALUES (?)`)
	if err != nil {
		fmt.Printf("Prepare sql err:%v\n", err)
		return 0
	}
	sqlResp, err := sqlPre.Exec(0)
	if err != nil {
		fmt.Printf("Exec sql err:%v\n", err)
		return 0
	}

	postId, err := sqlResp.LastInsertId()
	if err != nil {
		fmt.Printf("Get Last insert id err:%v\n", err)
		return 0
	}

	return postId
}

// 插入thread列表
func (d *DiscuzSql) InsertThread(article *Article) int64 {
	sqlPre, err := d.db.Prepare(`INSERT ` + d.dbPrefix + `forum_thread (fid, subject, author, authorid,
		dateline, lastpost, lastposter) VALUES (?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		fmt.Printf("Prepare sql err:%v\n", err)
		return 0
	}

	dateLine := time.Now().Unix()
	sqlResp, err := sqlPre.Exec(article.ClassId, article.Title, d.author, d.authorId, dateLine, dateLine, d.author)
	if err != nil {
		fmt.Printf("Exec sql err:%v\n", err)
		return 0
	}

	threadId, err := sqlResp.LastInsertId()
	if err != nil {
		fmt.Printf("Get Last insert id err:%v\n", err)
		return 0
	}

	return threadId
}

func (d *DiscuzSql) GenTags(article *Article, tid int64) {
	if article.Keywords == "" {
		return
	}

	keywords := strings.Split(article.Keywords, ",")
	formatKeywords := make([]string, 0)
	for _, keyword := range keywords {
		querySql := fmt.Sprintf(`SELECT tagid FROM `+d.dbPrefix+`common_tag WHERE tagname="%s" LIMIT 1`, keyword)
		rows, err := d.db.Query(querySql)
		if err != nil {
			fmt.Printf("query err:%v\n", err)
			return
		}
		defer rows.Close()

		var tagId int64
		if rows.Next() {
			rows.Scan(&tagId)
			formatKeywords = append(formatKeywords, fmt.Sprintf("%d,%s", tagId, keyword))
		} else {
			// 不存在tag需要先创建
			sqlPre, err := d.db.Prepare(`INSERT ` + d.dbPrefix + `common_tag (tagname) VALUES(?)`)
			if err != nil {
				fmt.Printf("Prepare GenTags sql err:%v\n", err)
				return
			}

			sqlResp, err := sqlPre.Exec(keyword)
			if err != nil {
				fmt.Printf("Exec GenTags sql err:%v\n", err)
				return
			}

			tagId, err = sqlResp.LastInsertId()
			if err != nil {
				fmt.Printf("GenTags get Last insert id err:%v\n", err)
				return
			}
			formatKeywords = append(formatKeywords, fmt.Sprintf("%d,%s", tagId, keyword))
		}

		// 更新tag和文章关系
		d.InsertTagItem(tagId, tid)
	}

	if len(formatKeywords) > 0 {
		article.Keywords = strings.Join(formatKeywords, "\t")
	}
}

func (d *DiscuzSql) InsertTagItem(tagId, tid int64) {
	sqlPre, err := d.db.Prepare(`INSERT ` + d.dbPrefix + `common_tagitem (tagid, itemid, idtype) VALUES(?, ?, "tid")`)
	if err != nil {
		fmt.Printf("Prepare InsertTagItem sql err:%v\n", err)
		return
	}
	sqlResp, err := sqlPre.Exec(tagId, tid)
	if err != nil {
		fmt.Printf("Exec InsertTagItem sql err:%v\n", err)
		return
	}

	rowsAffectedNum, err := sqlResp.RowsAffected()
	if err != nil || rowsAffectedNum == 0 {
		fmt.Printf("Rows affected err:%v\n", err)
		return
	}
}

func (d *DiscuzSql) InsertPost(article *Article, pid, tid int64) {
	sqlPre, err := d.db.Prepare(`INSERT ` + d.dbPrefix + `forum_post (pid, fid, tid, first,
		author, authorId, subject, dateline, message, usesig, smileyoff, position,
		tags) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		fmt.Printf("Prepare InsertPost sql err:%v\n", err)
		return
	}

	dateLine := time.Now().Unix()
	sqlResp, err := sqlPre.Exec(pid, article.ClassId, tid, 1, d.author, d.authorId, article.Title,
		dateLine, article.Body, 1, -1, 1, article.Keywords)
	if err != nil {
		fmt.Printf("Exec InsertPost sql err:%v\n", err)
		return
	}
	fmt.Printf("Exec InsertPost sql success:%v\n", sqlResp)
}

func (d *DiscuzSql) UpdateForum(article *Article, tid int64) {
	sqlPre, err := d.db.Prepare(`UPDATE ` + d.dbPrefix + `forum_forum SET threads=threads+1,
		posts=posts+1, todayposts=todayposts+1, lastpost=? WHERE fid=?`)
	if err != nil {
		fmt.Printf("Prepare UpdateForum sql err:%v\n", err)
		return
	}

	lastpost := fmt.Sprintf("%d", tid) + "\t" + article.Title + "\t" + fmt.Sprintf("%d", time.Now().Unix()) + "\t" + d.author
	sqlResp, err := sqlPre.Exec(lastpost, article.ClassId)
	if err != nil {
		fmt.Printf("Exec UpdateForum sql err:%v\n", err)
		return
	}
	fmt.Printf("Exec UpdateForum sql success:%v\n", sqlResp)
}

func (d *DiscuzSql) Close() {
	d.db.Close()
}
