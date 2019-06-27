package models

// User is a struct that contains all user info
type User struct {
	Id       uint64
	Email    string
	Password string
	Token    string
}

var id uint64

func NewUser(id uint64, email, password, token string) *User {

	return &User{
		Id:       id,
		Email:    email,
		Password: password,
		Token:    token,
	}
}
