package save

import (
	"errors"
	"net/http"
	"savemylink/models"
	save "savemylink/save/models"
	"strconv"

	readability "github.com/go-shiori/go-readability"
)

const (
	invalid_id = 0
)

type ArticleDB interface {
	SaveArticle(url, title, content, image string) (*save.Article, error)

	GetArticle(articleId uint64) (*save.Article, error)

	DeleteArticle(articleId uint64) error
}

type ArticleService interface {
	DownloadArticle(vc *models.ViewerContext, req *http.Request) (*save.Article, error)

	GetSavedArticle(vc *models.ViewerContext, req *http.Request) (*save.Article, error)

	DeleteArticle(vc *models.ViewerContext, req *http.Request) error
}

type Downloader func(url string) (*readability.Article, error)

type articleService struct {
	ArticleDB
	downloader Downloader
}

func NewArticleService(articleDB ArticleDB, downloader Downloader) ArticleService {
	return &articleService{
		ArticleDB:  articleDB,
		downloader: downloader,
	}
}

func (as *articleService) DownloadArticle(
	vc *models.ViewerContext,
	req *http.Request) (*save.Article, error) {
	if err := req.ParseForm(); err != nil {
		return nil, err
	}

	url := req.FormValue("url")
	article, err := as.downloader(url)
	if err != nil {
		return nil, err
	}

	return as.saveArticle(vc, url, article)
}

func (as *articleService) saveArticle(
	vc *models.ViewerContext,
	url string,
	article *readability.Article) (*save.Article, error) {

	return as.ArticleDB.SaveArticle(
		url,
		article.Title,
		article.Content,
		article.Image)
}

func (as *articleService) GetSavedArticle(
	vc *models.ViewerContext,
	req *http.Request) (*save.Article, error) {

	article_id, err := getArticleIdFromQuery(req)
	if err != nil {
		return nil, err
	}

	return as.ArticleDB.GetArticle(article_id)
}

func (as *articleService) DeleteArticle(
	vc *models.ViewerContext,
	req *http.Request) error {
	article_id, err := getArticleIdFromQuery(req)
	if err != nil {
		return err
	}

	return as.ArticleDB.DeleteArticle(article_id)
}

func getArticleIdFromQuery(req *http.Request) (uint64, error) {
	article_id, ok := req.URL.Query()["article_id"]
	if !ok || len(article_id[0]) < 1 {
		return invalid_id, errors.New("article_id is not provided")
	}

	id, err := strconv.Atoi(article_id[0])
	if err != nil {
		return invalid_id, err
	}

	return uint64(id), nil
}
