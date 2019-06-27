package save

import (
	"html"
	"net/http"
	"net/url"
	"strings"

	readability "github.com/go-shiori/go-readability"
)

func getContent(urlString string) (*readability.Article, error) {
	if correct, err := isCorrect(urlString); !correct {
		return nil, err
	}

	return downloadContent(cleanUrl(urlString))
}

func isCorrect(urlString string) (bool, error) {
	_, err := url.ParseRequestURI(urlString)
	return err == nil, err
}

func cleanUrl(urlString string) string {
	htmlDecoded := htmlEntityDecode(urlString)

	pos := stringPos(htmlDecoded, "&utm_source=", 0)
	cleanUrl := htmlDecoded
	if pos != -1 {
		cleanUrl = htmlDecoded[:pos]
	}

	pos = stringPos(cleanUrl, "?utm_source=", 0)
	if pos != -1 {
		cleanUrl = cleanUrl[:pos]
	}

	pos = stringPos(cleanUrl, "#xtor=RSS-", 0)
	if pos != -1 {
		cleanUrl = cleanUrl[:pos]
	}

	return cleanUrl
}

func htmlEntityDecode(str string) string {
	return html.UnescapeString(str)
}

func stringPos(haystack, needle string, offset int) int {
	length := len(haystack)
	if length == 0 || offset > length || -offset > length {
		return -1
	}

	if offset < 0 {
		offset += length
	}

	pos := strings.Index(haystack[offset:], needle)
	if pos == -1 {
		return -1
	}

	return pos + offset
}

func downloadContent(urlString string) (*readability.Article, error) {
	response, err := http.Get(urlString)
	defer response.Body.Close()

	if err != nil {
		return nil, err
	}

	article, err := readability.FromReader(response.Body, urlString)
	if err != nil {
		return nil, err
	}

	return &article, nil
}
