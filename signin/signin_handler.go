package signin

import (
	"encoding/json"
	"net/http"
	"savemylink/services"
)

type SignupHandler struct {
	UserService services.UserService
}

func NewSignupHandler(userDB services.UserDB) *SignupHandler {
	return &SignupHandler{
		UserService: services.NewUserService(userDB),
	}
}

func (sh *SignupHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		sh.signUp(res, req)
		return
	}
}

func (sh *SignupHandler) signUp(res http.ResponseWriter, req *http.Request) {
	user, err := sh.UserService.CreateNewUser(res, req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(user)
}
