package utils

import (
	"fmt"
	"testing"
)

func TestGenerateUsername(t *testing.T) {
	tests := []struct {
		email          string
		expectedResult string
		expectedError  error
	}{
		{"john.doe@example.com", "gopherjohndoe", nil},
		{"jane.smith@example.com", "gopherjanesmith", nil},
		{"first.last@example.com", "gopherfirstlast", nil},
		{"username@example.com", "gopherusername", nil},
		{"invalidemailformat", "", fmt.Errorf("invalid email format")},
		{"singlepart@com", "gophersinglepart", nil},
	}

	for _, test := range tests {
		username, err := GenerateUsername(test.email)

		if username != test.expectedResult {
			t.Errorf("For %s, expected username %s, but got %s", test.email, test.expectedResult, username)
		}

		if (err == nil && test.expectedError != nil) || (err != nil && test.expectedError == nil) || (err != nil && err.Error() != test.expectedError.Error()) {
			t.Errorf("For %s, expected error %v, but got %v", test.email, test.expectedError, err)
		}
	}
}
