package database

import (
	"strings"
	"testing"
)

var email = "me@email.com"
var password = "password"

func TestCreateNewUser(t *testing.T) {
	db := NewMemoryDb()
	user, error := db.CreateUser(email, password)

	printErrorsIfAny(t, error)

	if strings.Compare(email, user.Email) != 0 {
		t.Errorf("User with email %s is created. Wanted:%s", user.Email, email)
	}

	if strings.Compare(password, user.Password) != 0 {
		t.Errorf("User with email %s is created. Wanted:%s", user.Email, email)
	}
}

func TestDuplicateUserCreation(t *testing.T) {
	db := NewMemoryDb()
	db.CreateUser(email, password)

	// Create a new user with the same email and password again
	user, err := db.CreateUser(email, password)
	if user != nil && err == nil {
		t.Errorf("User %d is created again", user)
	}
}

func TestGetUserById(t *testing.T) {
	db := NewMemoryDb()
	user, _ := db.CreateUser(email, password)

	foundUser, err := db.GetById(user.Id)
	printErrorsIfAny(t, err)
	if strings.Compare(email, foundUser.Email) != 0 {
		t.Errorf("User with email %s wasn't found", foundUser.Email)
	}
}

func TestGetUserByEmail(t *testing.T) {
	db := NewMemoryDb()
	user, _ := db.CreateUser(email, password)

	foundUser, err := db.GetByEmail(user.Email)
	printErrorsIfAny(t, err)
	if strings.Compare(email, foundUser.Email) != 0 {
		t.Errorf("User with email %s wasn't found", foundUser.Email)
	}
}

func printErrorsIfAny(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Error %d found", err)
	}
}
