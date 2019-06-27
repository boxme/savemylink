package save

import "testing"

func TestUrlCheck(t *testing.T) {
	tests := []struct {
		url             string
		expected_result bool
	}{
		{"https://medium.com/@saiyerram/go-interfaces-pointers-4d1d98d5c9c6", true},
		{"google.com", false},
	}

	for _, tt := range tests {
		output, _ := isCorrect(tt.url)

		if output != tt.expected_result {
			t.Errorf("%s is %d, but got:%d", output, tt.expected_result, output)
		}
	}
}
