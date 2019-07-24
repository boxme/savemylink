package auth

import (
	"net/http"
	"savemylink/database"
	"savemylink/login"
	"savemylink/logout"
	"savemylink/models"
	"savemylink/save"
	"savemylink/services"
	"savemylink/signin"
	"savemylink/util"
)

// Handler is a struct that contains all other handlers that will handle
// subsequent requests
type Handler struct {
	UserService     services.UserService
	SignupHandler   *signin.SignupHandler
	LoginHandler    *login.LoginHandler
	LogoutHandler   *logout.LogoutHandler
	requestHandlers []RequestHandler
}

// RequestHandler is an interface that handles the request at the highest level
type RequestHandler interface {
	// Return true if request is handled.
	ServeHTTP(vc *models.ViewerContext, res http.ResponseWriter, req *http.Request) bool
}

// NewAuthHandler creates a new instance of authentication handler
func NewAuthHandler(db database.Db) *Handler {
	reqHandlers := []RequestHandler{save.NewSaveHandler(db.GetArticleDb())}
	userDb := db.GetUserDb()
	return &Handler{
		UserService:     services.NewUserService(userDb),
		SignupHandler:   signin.NewSignupHandler(userDb),
		LoginHandler:    login.NewSignupHandler(userDb),
		LogoutHandler:   logout.NewLogoutHandler(userDb),
		requestHandlers: reqHandlers,
	}
}

func (ah *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = util.ShiftPath(req.URL.Path)
	if head == "logout" {
		ah.LoginHandler.ServeHTTP(res, req)
		return
	}

	viewerContext := ah.authenticate(res, req)
	if viewerContext != nil {
		ah.route(viewerContext, res, req)
		return
	}

	if head == "signup" {
		ah.SignupHandler.ServeHTTP(res, req)
		return
	}
	if head == "login" {
		ah.LoginHandler.ServeHTTP(res, req)
		return
	}
}

func (ah *Handler) authenticate(
	res http.ResponseWriter,
	req *http.Request) *models.ViewerContext {

	cookie, err := req.Cookie("user_token")
	if err != nil {
		return nil
	}

	userToken := cookie.Value
	user, err := ah.UserService.GetByUserToken(userToken)
	if err != nil {
		return nil
	}

	return models.NewViewerContext(user.Id, user.Token)
}

func (ah *Handler) route(
	vc *models.ViewerContext, res http.ResponseWriter, req *http.Request) {
	// Do other routing
	for _, reqHandler := range ah.requestHandlers {
		if reqHandler.ServeHTTP(vc, res, req) {
			return
		}
	}
}
