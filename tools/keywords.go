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

func NewKeywordsHandler(dict string) *KeywordsHandler {
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
	for i, keyword := range segStrs {
		if len(keyword) < MinKeywordLen {
			continue
		}

		if !isNouns(keyword) {
			continue
		}

		if i > 0 && isAdjectiveWord(segStrs[i-1]) {
			keywords = append(keywords, removeTail(segStrs[i-1])+removeTail(keyword))
		} else {
			keywords = append(keywords, removeTail(keyword))
		}

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

func removeTail(keyword string) string {
	return keyword[:strings.Index(keyword, "/")]
}
