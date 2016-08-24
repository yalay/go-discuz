package ecms

import (
	"encoding/json"
	"fmt"
	"regexp"
)

type Article struct {
	ClassId  int
	Title    string
	Keywords string
	Body     string
}

func NewArticle(classId int, title, keywords, body string) *Article {
	if title == "" {
		return nil
	}
	return &Article{
		ClassId:  classId,
		Title:    title,
		Keywords: keywords,
		Body:     body,
	}
}

func (article *Article) String() string {
	if article == nil {
		return ""
	}

	jsonArticle, err := json.Marshal(article)
	if err != nil {
		fmt.Printf("json Marshal err:%v", err)
		return ""
	}

	return string(jsonArticle)
}

func (article *Article) FormatBodyForDiscuz() {
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
			return "[img]" + newImg + "[/img]\n"
		})
}
