package save

import (
	"net/http"
	"savemylink/models"
	save "savemylink/save/models"
	"savemylink/util"
)

// Handler is a Struct that contains an instance of ArticleService
type Handler struct {
	articleService ArticleService
}

// NewSaveHandler creates a new save request handler
func NewSaveHandler(articleDb ArticleDB) *Handler {
	return &Handler{
		articleService: NewArticleService(articleDb, getContent),
	}
}

func (sh *Handler) ServeHTTP(
	vc *models.ViewerContext, res http.ResponseWriter, req *http.Request) bool {

	var head string
	head, req.URL.Path = util.ShiftPath(req.URL.Path)

	if head != "save" {
		return false
	}

	switch req.Method {
	case "POST":
		sh.saveArticle(vc, res, req)
		break
	case "GET":
		sh.getArticle(vc, res, req)
		break
	case "DELETE":
		sh.deleteArticle(vc, res, req)
		break
	}

	return true
}

func (sh *Handler) saveArticle(
	vc *models.ViewerContext, res http.ResponseWriter, req *http.Request) {
	_, err := sh.articleService.DownloadArticle(vc, req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
}

func (sh *Handler) getArticle(
	vc *models.ViewerContext,
	res http.ResponseWriter,
	req *http.Request) *save.Article {
	article, err := sh.articleService.GetSavedArticle(vc, req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return nil
	}

	return article
}

func (sh *Handler) deleteArticle(
	vc *models.ViewerContext,
	res http.ResponseWriter,
	req *http.Request) {
	err := sh.articleService.DeleteArticle(vc, req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}
}
