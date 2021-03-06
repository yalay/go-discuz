package discuz

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
	"tools"
)

type Article struct {
	ClassId  int
	Title    string
	Cover    string
	Keywords string
	Body     string
}

func PublishArticleFromFile(config *tools.Config) {
	if config == nil || config.DataFile == "" {
		fmt.Println("config is empty.\n")
		return
	}

	articleFile, err := os.Open(config.DataFile)
	if err != nil {
		fmt.Printf("Open file err:%v\n", err)
		return
	}
	defer articleFile.Close()

	discuzSql := newDiscuzSql(config.DiscuzSql)
	if discuzSql == nil {
		return
	}
	defer discuzSql.Close()

	buff := bufio.NewReader(articleFile)
	for {
		jsonArticle, err := buff.ReadString('\n')
		if err != nil && len(jsonArticle) == 0 {
			break
		}

		jsonArticle = strings.TrimSpace(jsonArticle)
		article := newArticleFromJson(jsonArticle)
		if article == nil {
			continue
		}

		if config.EnableFidMapping {
			article.mappingClassId(config)
		}

		if config.EnableFormatImg {
			article.formatImgLabel(config)
		}

		if config.EnableGenCover {
			article.genCover(config)
		}

		if config.EnableGenKeyword {
			handler := tools.GetKeywordsHandler(config.Dict)
			article.genKeyWords(handler)
		}

		article.publish(discuzSql)
	}
}

func newArticleFromJson(jsonArticle string) *Article {
	if jsonArticle == "" {
		return nil
	}

	var newArticle = &Article{}
	err := json.Unmarshal([]byte(jsonArticle), newArticle)
	if err != nil {
		fmt.Printf("json Unmarshal err:%v\n", err)
		return nil
	}
	return newArticle
}

// 按照discuz的格式格式化内容
func (article *Article) formatImgLabel(config *tools.Config) {
	if article == nil || article.Body == "" {
		return
	}

	// 调整img标签内容
	reg := regexp.MustCompile(`<img.+?\/>`)
	article.Body = reg.ReplaceAllStringFunc(article.Body,
		func(oriImg string) string {
			imgReg := regexp.MustCompile(`http:\S+\.(?i:jpg|jpeg|gif|png|webp)`)
			newImg := imgReg.FindString(oriImg)
			if newImg == "" {
				return ""
			}
			if !config.IsImgWhiteList(newImg) {
				return ""
			}

			return "[img]" + newImg + "[/img]\n"
		})
}

// 从标题里面提取关键词
func (article *Article) genKeyWords(handler *tools.KeywordsHandler) {
	if handler == nil {
		return
	}

	keywords := handler.GetKeywords(article.Title)
	if len(keywords) >= tools.MinKeywordLen {
		article.Keywords = keywords
	}
}

func (article *Article) mappingClassId(config *tools.Config) {
	fid := config.GetMappingFid(article.ClassId)
	if fid != 0 {
		article.ClassId = fid
	}
}

func (article *Article) genCover(config *tools.Config) {
	if article.Cover != "" {
		return
	}

	coverReg := regexp.MustCompile(`http:\S+\.(?i:jpg|jpeg|gif|png|webp)`)
	coverImg := coverReg.FindString(article.Body)
	if coverImg != "" {
		article.Cover = coverImg + config.ThumbParam
	}
}

func (article *Article) publish(discuzSql *DiscuzSql) {
	if discuzSql.CheckTitleExist(article) {
		return
	}

	pid := discuzSql.GetPostId()
	if pid == 0 {
		fmt.Printf("get pid err. Article:+%v\n", article)
		return
	}

	tid := discuzSql.InsertThread(article)
	if tid == 0 {
		fmt.Printf("insert thread err. Article:+%v\n", article)
		return
	}

	discuzSql.InsertCover(article, tid)
	discuzSql.GenTags(article, tid)
	discuzSql.InsertPost(article, pid, tid)
	time.Sleep(100 * time.Millisecond)
}
