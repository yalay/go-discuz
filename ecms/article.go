package ecms

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

type Article struct {
	Id       int
	ClassId  int
	Title    string
	Cover    string
	Keywords string
	Body     string
}

func NewArticle(id, classId int, title, cover, keywords, body string) *Article {
	if title == "" {
		return nil
	}
	return &Article{
		Id:       id,
		ClassId:  classId,
		Title:    title,
		Cover:    cover,
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

func (article *Article) Dump(filename string) {
	articleFile, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open file err:%v", err)
		return
	}
	defer articleFile.Close()

	if _, err = articleFile.WriteString(article.String()); err != nil {
		fmt.Printf("write file err:%v", err)
		return
	}
	articleFile.WriteString("\n")
}
