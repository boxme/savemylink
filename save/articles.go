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
	invalidID = 0
)

// ArticleDB is an interface for articles persistence storage
type ArticleDB interface {
	SaveArticle(url, title, content, image string) (*save.Article, error)

	GetArticle(articleID uint64) (*save.Article, error)

	DeleteArticle(articleID uint64) error
}

// ArticleService is an interface for handling articles
type ArticleService interface {
	DownloadArticle(vc *models.ViewerContext, req *http.Request) (*save.Article, error)

	GetSavedArticle(vc *models.ViewerContext, req *http.Request) (*save.Article, error)

	DeleteArticle(vc *models.ViewerContext, req *http.Request) error
}

// Downloader is a function type for downloading an article from its url
type Downloader func(url string) (*readability.Article, error)

type articleService struct {
	ArticleDB
	downloader Downloader
}

// NewArticleService creates a new instance of ArticleService
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

	articleID, err := getArticleIDFromQuery(req)
	if err != nil {
		return nil, err
	}

	return as.ArticleDB.GetArticle(articleID)
}

func (as *articleService) DeleteArticle(
	vc *models.ViewerContext,
	req *http.Request) error {
	articleID, err := getArticleIDFromQuery(req)
	if err != nil {
		return err
	}

	return as.ArticleDB.DeleteArticle(articleID)
}

func getArticleIDFromQuery(req *http.Request) (uint64, error) {
	articleID, ok := req.URL.Query()["articleID"]
	if !ok || len(articleID[0]) < 1 {
		return invalidID, errors.New("articleID is not provided")
	}

	id, err := strconv.Atoi(articleID[0])
	if err != nil {
		return invalidID, err
	}

	return uint64(id), nil
}
