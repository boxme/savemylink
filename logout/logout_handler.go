package logout

import (
	"net/http"
	"savemylink/services"
)

type LogoutHandler struct {
	userService services.UserService
}

func NewLogoutHandler(userDB services.UserDB) *LogoutHandler {
	return &LogoutHandler{
		userService: services.NewUserService(userDB),
	}
}

func (lh *LogoutHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		lh.logout(res, req)
		return
	}
}

func (lh *LogoutHandler) logout(res http.ResponseWriter, req *http.Request) {
	err := lh.userService.Logout(res, req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
}
