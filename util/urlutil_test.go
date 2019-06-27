package util

import (
	"strings"
	"testing"
)

func TestShiftPath(t *testing.T) {
	tests := []struct {
		input         string
		expected_head string
		expected_tail string
	}{
		{"www.savemylink.com", "www.savemylink.com", "/"},
		{"www.savemylink.com/user", "www.savemylink.com", "/user"},
		{"www.savemylink.com/user/profile", "www.savemylink.com", "/user/profile"},
		{"user/profile?id=3", "user", "/profile?id=3"},
		{"/profile?id=3", "profile?id=3", "/"},
	}

	for _, tt := range tests {
		output_head, output_tail := ShiftPath(tt.input)
		isEqual := strings.Compare(output_head, tt.expected_head) == 0
		outputErrors(t, isEqual, "head", output_head, tt.expected_head)

		isEqual = strings.Compare(output_tail, tt.expected_tail) == 0
		outputErrors(t, isEqual, "tail", output_tail, tt.expected_tail)
	}
}

func outputErrors(t *testing.T, isEqual bool, output_type, output, expected string) {
	if !isEqual {
		t.Errorf("%s:%s, require:%s", output_type, output, expected)
	}
}
