package spider

import (
	"bytes"
	"net/url"

	"golang.org/x/net/html"
)

const (
	aTag     = "a"
	baseTag  = "base"
	hrefAttr = "href"
)

func findAHref(tokenizer *html.Tokenizer, urls *[]string, base string) {
	hasMore := true
	var key, value []byte
	for hasMore {
		key, value, hasMore = tokenizer.TagAttr()
		keyS, valueS := string(key), string(value)
		if keyS == hrefAttr {
			urlParsed, _ := url.Parse(valueS)
			if base != "" && urlParsed.Host == "" {
				valueS = base + valueS
			}
			*urls = append(*urls, valueS)
		}
	}
}

func findBaseHref(tokenizer *html.Tokenizer) string {
	hasMore := true
	var key, value []byte
	for hasMore {
		key, value, hasMore = tokenizer.TagAttr()
		keyS, valueS := string(key), string(value)
		if keyS == hrefAttr {
			return valueS
		}
	}
	return ""
}

func parseUrls(htmlBytes []byte) (urls []string) {
	tokenizer := html.NewTokenizer(bytes.NewBuffer(htmlBytes))
	urls = []string{}
	base := ""
	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			return
		case html.StartTagToken:
			tn, hasAttr := tokenizer.TagName()
			tnS := string(tn)
			if tnS == baseTag && hasAttr {
				base = findBaseHref(tokenizer)
			}
			if tnS == aTag && hasAttr {
				findAHref(tokenizer, &urls, base)
			}
		}
	}
	return
}
