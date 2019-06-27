package signin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"savemylink/database"
	"savemylink/models"
	"strings"
	"testing"
)

func TestSignin(t *testing.T) {
	email := "test@email"
	password := "password"

	form := url.Values{}
	form.Add("email", email)
	form.Add("password", password)

	req, err := http.NewRequest(
		"POST",
		"/signin",
		strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	signup_handler := NewSignupHandler(database.NewMemoryDb().GetUserDb())
	signup_handler.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusCreated {
		t.Errorf(
			"Return wrong status code. got %v wanted %v",
			responseRecorder.Code,
			http.StatusCreated)
	}

	result := responseRecorder.Result()
	defer result.Body.Close()
	decoder := json.NewDecoder(result.Body)

	var user models.User
	decoder.Decode(&user)

	if user.Email != email || user.Password != password {
		t.Errorf("handler returned %v", user)
	}

	expectedCType := "application/json"
	cType := responseRecorder.Header().Get("Content-Type")
	if cType != expectedCType {
		t.Errorf(
			"content type header does not match: got %v want %v",
			cType,
			expectedCType)
	}
}
