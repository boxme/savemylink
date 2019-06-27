package save

import "testing"

func TestUrlCheck(t *testing.T) {
	tests := []struct {
		url             string
		expectedResult bool
	}{
		{"https://medium.com/@saiyerram/go-interfaces-pointers-4d1d98d5c9c6", true},
		{"google.com", false},
	}

	for _, tt := range tests {
		output, _ := isCorrect(tt.url)

		if output != tt.expectedResult {
			t.Errorf("%s is %t, but got %t", tt.url, tt.expectedResult, output)
		}
	}
}
