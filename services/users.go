package services

import (
	"errors"
	"net/http"
	"savemylink/models"
	"savemylink/util"
)

type UserDB interface {
	CreateUser(email, password, userToken string) (*models.User, error)
	LoginUser(email, userToken string) (*models.User, error)
	GetById(id uint64) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByUserToken(userToken string) (*models.User, error)
	LogoutUser(userToken string) error
}

type UserService interface {
	CreateNewUser(res http.ResponseWriter, req *http.Request) (*models.User, error)
	Login(res http.ResponseWriter, req *http.Request) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserById(id uint64) (*models.User, error)
	GetByUserToken(userToken string) (*models.User, error)
	Logout(res http.ResponseWriter, req *http.Request) error
}

type userService struct {
	UserDB
}

func NewUserService(userDB UserDB) UserService {
	return &userService{
		UserDB: userDB,
	}
}

func (us *userService) CreateNewUser(res http.ResponseWriter, req *http.Request) (*models.User, error) {
	if err := req.ParseForm(); err != nil {
		return nil, err
	}

	email := req.FormValue("email")
	if oldUser, err := us.UserDB.GetByEmail(email); oldUser != nil {
		return nil, err
	}

	if email == "" {
		return nil, errors.New("no email is given")
	}

	userToken, err := util.GenerateUserToken()
	if err != nil {
		return nil, err
	}

	password := req.FormValue("password")
	if password == "" {
		return nil, errors.New("no password is given")
	}

	newUser, err := us.UserDB.CreateUser(email, password, userToken)
	setUserTokenToCookie(newUser, userToken, res)

	return newUser, err
}

func (us *userService) Login(
	res http.ResponseWriter,
	req *http.Request) (*models.User, error) {
	// TODO: Validate password
	if err := req.ParseForm(); err != nil {
		return nil, err
	}

	userToken, err := util.GenerateUserToken()
	if err != nil {
		return nil, err
	}

	user, err := us.UserDB.LoginUser(req.FormValue("email"), userToken)
	if err != nil {
		return nil, err
	}

	setUserTokenToCookie(user, userToken, res)

	return user, nil
}

func (us *userService) Logout(
	res http.ResponseWriter,
	req *http.Request) error {

	cookie, err := req.Cookie("user_token")
	if err != nil {
		return err
	}

	if err := us.LogoutUser(cookie.Value); err != nil {
		return err
	}

	removeUserTokenFromCookie(res)

	return nil
}

func (us *userService) GetUserByEmail(email string) (*models.User, error) {
	// TODO: Validate email
	return us.UserDB.GetByEmail(email)
}

func (us *userService) GetUserById(id uint64) (*models.User, error) {
	// TODO: Validate
	return us.UserDB.GetById(id)
}

func (us *userService) GetByUserToken(userToken string) (*models.User, error) {
	return us.GetByUserToken(userToken)
}

func (us *userService) LogoutUser(userToken string) error {
	return us.LogoutUser(userToken)
}

func setUserTokenToCookie(user *models.User, userToken string, res http.ResponseWriter) {
	if user != nil {
		cookie := http.Cookie{
			Name:     "user_token",
			Value:    userToken,
			HttpOnly: true,
		}
		http.SetCookie(res, &cookie)
	}
}

func removeUserTokenFromCookie(res http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     "user_token",
		Value:    "",
		MaxAge:   0,
		HttpOnly: true,
	}
	http.SetCookie(res, &cookie)
}
