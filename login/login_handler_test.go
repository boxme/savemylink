package login

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"savemylink/database"
	"savemylink/models"
	"savemylink/services"
	"strings"
	"testing"
)

func TestLogin(t *testing.T) {
	email := "email"
	password := "password"
	token := "token"
	userDB := database.NewMemoryDb().GetUserDb()
	createUser(email, password, token, userDB)

	form := url.Values{}
	form.Add("email", email)
	form.Add("password", password)
	req, err := http.NewRequest(
		"POST",
		"/login",
		strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	loginHandler := NewSignupHandler(userDB)
	loginHandler.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf(
			"Return wrong status code. got %v wanted %v",
			responseRecorder.Code,
			http.StatusOK)
	}

	expectedCType := "application/json"
	cType := responseRecorder.Header().Get("Content-Type")
	if cType != expectedCType {
		t.Errorf(
			"content type header does not match: got %v want %v",
			cType,
			expectedCType)
	}

	result := responseRecorder.Result()
	defer result.Body.Close()
	decoder := json.NewDecoder(result.Body)

	var user models.User
	decoder.Decode(&user)
	if user.Email != email || user.Password != password {
		t.Errorf("handler returned %v", user)
	}

	cookie := result.Cookies()[0]
	if cookie.Name != "user_token" && cookie.Value != user.Token {
		t.Errorf(
			"Return wrong user token: got %v ant %v",
			cookie.Value,
			user.Token)
	}
}

func createUser(email, password, token string, userDB services.UserDB) {
	userDB.CreateUser(email, password, token)
}
