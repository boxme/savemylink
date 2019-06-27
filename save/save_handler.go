package save

import (
	"net/http"
	"savemylink/models"
	save "savemylink/save/models"
	"savemylink/util"
)

type SaveHandler struct {
	articleService ArticleService
}

func NewSaveHandler(articleDb ArticleDB) *SaveHandler {
	return &SaveHandler{
		articleService: NewArticleService(articleDb, getContent),
	}
}

func (sh *SaveHandler) ServeHTTP(
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

func (sh *SaveHandler) saveArticle(
	vc *models.ViewerContext, res http.ResponseWriter, req *http.Request) {
	_, err := sh.articleService.DownloadArticle(vc, req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
}

func (sh *SaveHandler) getArticle(
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

func (sh *SaveHandler) deleteArticle(
	vc *models.ViewerContext,
	res http.ResponseWriter,
	req *http.Request) {
	err := sh.articleService.DeleteArticle(vc, req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}
}
