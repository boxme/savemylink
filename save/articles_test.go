package save

import (
	"errors"
	"net/http"
	"net/url"
	"savemylink/models"
	model "savemylink/save/models"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	readability "github.com/go-shiori/go-readability"
)

const (
	content     = "content"
	downloadURL = "test_url"
)

func testDownload(t *testing.T) {
	article, _ := mockDownloadContent(t, downloadURL)

	if article.Url != downloadURL {
		t.Errorf(
			"Return wrong url: got %v want %v",
			article.Url,
			downloadURL)
	}

	if article.Content != content {
		t.Errorf(
			"Return wrong content: got %v want %v",
			article.Content,
			content)
	}
}

func testGetSavedArticle(t *testing.T) {
	article, articleService := mockDownloadContent(t, downloadURL)

	req, err := http.NewRequest("GET", "/save/", nil)
	if err != nil {
		t.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("article_id", strconv.FormatUint(article.Id, 10))
	req.URL.RawQuery = q.Encode()

	persistedArticle, err :=
		articleService.GetSavedArticle(
			models.NewViewerContext(1, "token"), req)

	if err != nil {
		t.Fatal(err)
	}

	if article.Content != persistedArticle.Content {
		t.Errorf(
			"Return wrong content: got %v want %v",
			article.Content,
			persistedArticle.Content)
	}
}

func testDeleteArticle(t *testing.T) {
	article, articleService := mockDownloadContent(t, downloadURL)

	req, err := http.NewRequest("DELETE", "/save/", nil)
	if err != nil {
		t.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("article_id", strconv.FormatUint(article.Id, 10))
	req.URL.RawQuery = q.Encode()

	err =
		articleService.DeleteArticle(
			models.NewViewerContext(1, "token"), req)

	if err != nil {
		t.Fatal(err)
	}
}

func mockDownloadContent(
	t *testing.T,
	downloadURL string) (*model.Article, ArticleService) {
	form := url.Values{}
	form.Add("url", downloadURL)
	req, err := http.NewRequest(
		"POST",
		"/save",
		strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	fakeDb := fakeArticlesDb()
	articleService :=
		NewArticleService(fakeDb, mockDownloader)

	article, err := articleService.DownloadArticle(nil, req)
	if err != nil {
		t.Fatal(err)
	}

	return article, articleService
}

func mockDownloader(url string) (*readability.Article, error) {
	article, err := readability.FromReader(strings.NewReader(content), url)
	return &article, err
}

func fakeArticlesDb() *articleMemoryDb {
	return &articleMemoryDb{
		id:           0,
		articleTable: make(map[string]*model.Article),
	}
}

type articleMemoryDb struct {
	id           uint64
	articleTable map[string]*model.Article
}

func (articleDb *articleMemoryDb) SaveArticle(
	url, title, content, image string) (*model.Article, error) {
	if _, ok := articleDb.articleTable[url]; ok {
		return nil, errors.New("article is already saved")
	}

	article :=
		model.NewArticle(
			atomic.AddUint64(&articleDb.id, 1), url, title, content, image)
	articleDb.articleTable[url] = article

	return article, nil
}

func (articleDb *articleMemoryDb) GetArticle(id uint64) (*model.Article, error) {
	for _, article := range articleDb.articleTable {
		if article.Id == id {
			return article, nil
		}
	}

	return nil, errors.New("Article not found")
}

func (articleDb *articleMemoryDb) DeleteArticle(id uint64) error {
	article, err := articleDb.GetArticle(id)
	if err == nil {
		delete(articleDb.articleTable, article.Url)
		return nil
	}

	return err
}
