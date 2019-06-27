package login

import (
	"encoding/json"
	"net/http"
	"savemylink/services"
)

type LoginHandler struct {
	userService services.UserService
}

func NewSignupHandler(userDB services.UserDB) *LoginHandler {
	return &LoginHandler{
		userService: services.NewUserService(userDB),
	}
}

func (lh *LoginHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		lh.login(res, req)
		return
	}
}

func (lh *LoginHandler) login(res http.ResponseWriter, req *http.Request) {
	user, err := lh.userService.Login(res, req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(user)
}
