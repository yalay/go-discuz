package tools

import (
	"github.com/huichen/sego"
	"strings"
)

const (
	MinKeywordLen = 4
	MaxKeywordNum = 5
)

type KeywordsHandler struct {
	seg sego.Segmenter
}

var keywordHandler *KeywordsHandler

func GetKeywordsHandler(dict string) *KeywordsHandler {
	if keywordHandler == nil {
		keywordHandler = newKeywordsHandler(dict)
	}
	return keywordHandler
}

func newKeywordsHandler(dict string) *KeywordsHandler {
	if dict == "" {
		return nil
	}

	var seg sego.Segmenter
	seg.LoadDictionary(dict)
	segDict := seg.Dictionary()
	if segDict == nil || segDict.NumTokens() == 0 {
		return nil
	}

	return &KeywordsHandler{
		seg: seg,
	}
}

func (handler *KeywordsHandler) GetKeywords(text string) string {
	segments := handler.seg.Segment([]byte(text))
	segStr := sego.SegmentsToString(segments, false)
	if segStr == "" {
		return ""
	}

	keywords := make([]string, 0)
	segStrs := strings.Fields(segStr)
	for i, keywordAttr := range segStrs {
		if !isNouns(keywordAttr) {
			continue
		}

		keyword := removeTail(keywordAttr)
		if i > 0 {
			lastkeywordAttr := segStrs[i-1]
			if isAdjectiveWord(lastkeywordAttr) || isVerb(lastkeywordAttr) {
				keyword = removeTail(lastkeywordAttr) + keyword
			}
		}

		if len(keyword) < MinKeywordLen {
			continue
		}

		keywords = append(keywords, keyword)
		if len(keywords) > MaxKeywordNum {
			break
		}
	}

	return strings.Join(keywords, ",")
}

func isNouns(keyword string) bool {
	if strings.HasSuffix(keyword, "/n") ||
		strings.HasSuffix(keyword, "/nr") ||
		strings.HasSuffix(keyword, "/ns") ||
		strings.HasSuffix(keyword, "/nt") ||
		strings.HasSuffix(keyword, "/nz") ||
		strings.HasSuffix(keyword, "/ng") {
		return true
	}
	return false
}

func isAdjectiveWord(keyword string) bool {
	if strings.HasSuffix(keyword, "/a") ||
		strings.HasSuffix(keyword, "/ad") ||
		strings.HasSuffix(keyword, "/an") ||
		strings.HasSuffix(keyword, "/ag") ||
		strings.HasSuffix(keyword, "/al") {
		return true
	}
	return false
}

func isVerb(keyword string) bool {
	if strings.HasSuffix(keyword, "/v") ||
		strings.HasSuffix(keyword, "/vd") ||
		strings.HasSuffix(keyword, "/vn") {
		return true
	}
	return false
}

func removeTail(keyword string) string {
	tailIdx := strings.Index(keyword, "/")
	if tailIdx > 0 {
		return keyword[:strings.Index(keyword, "/")]
	}
	return keyword
}
