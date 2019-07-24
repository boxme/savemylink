package util

import (
	"strings"
	"testing"
)

func TestShiftPath(t *testing.T) {
	tests := []struct {
		input        string
		expectedHead string
		expectedTail string
	}{
		{"www.savemylink.com", "www.savemylink.com", "/"},
		{"www.savemylink.com/user", "www.savemylink.com", "/user"},
		{"www.savemylink.com/user/profile", "www.savemylink.com", "/user/profile"},
		{"user/profile?id=3", "user", "/profile?id=3"},
		{"/profile?id=3", "profile?id=3", "/"},
	}

	for _, tt := range tests {
		outputHead, outputTail := ShiftPath(tt.input)
		isEqual := strings.Compare(outputHead, tt.expectedHead) == 0
		outputErrors(t, isEqual, "head", outputHead, tt.expectedHead)

		isEqual = strings.Compare(outputTail, tt.expectedTail) == 0
		outputErrors(t, isEqual, "tail", outputTail, tt.expectedTail)
	}
}

func outputErrors(t *testing.T, isEqual bool, outputType, output, expected string) {
	if !isEqual {
		t.Errorf("%s:%s, require:%s", outputType, output, expected)
	}
}
