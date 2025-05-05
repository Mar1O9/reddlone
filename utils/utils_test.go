package utils_test

import (
	"strings"
	"testing"
	"time"

	"github.com/Mar1O9/reddlone/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name       string
		username   string
		wantErr    bool
		errMessage string
	}{
		{"valid username", "validUsername", false, ""},
		{"invalid username - too short", "a", true, "invalid Username"},
		{"invalid username - too long", "KSJDKFLAJDFLKJAKLDFJKSDJFKSJDKJAKLJKASJDFKJASKLDFJASKDJFKASJDFKJASDKFJ", true, "invalid Username"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := utils.ValidateUsername(test.username)
			if err != nil && !test.wantErr {
				t.Errorf("ValidateUsername() unexpected error = %v", err)
			} else if err == nil && test.wantErr {
				t.Errorf("ValidateUsername() expected error, got nil")
			}
		})
	}
}

func TestValidateName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"empty string", "", true},
		{"less than 3 characters", "ab", true},
		{"more than 32 characters", strings.Repeat("a", 33), true},
		{"invalid characters", "abc!", false}, // Note: ! is a valid character
		{"invalid characters", "abc@", false}, // Note: @ is not a valid character
		{"valid characters", "abcdef", false},
		{"special characters", "abc!@#$%^&*()_+={}[", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := utils.ValidateName(test.input)
			if (err != nil) != test.wantErr {
				t.Errorf("ValidateName() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"valid email", "test@example.com", false},
		{"invalid email (missing @)", "testexample.com", true},
		{"invalid email (missing domain)", "test@", true},
		{"invalid email (invalid characters)", "test@example!com", true},
		{"empty string", "", true},
		{"nil string", "", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := utils.ValidateEmail(test.email)
			if (err != nil) != test.wantErr {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name           string
		password       string
		shouldValidate bool
	}{
		{"uppercase and digit", "P@ssw0rd", true},
		{"only lowercase", "password", false},
		{"only digits", "12345678", false},
		{"special characters only", "!@#$%^&*()", false},
		{"short password", "P@ss", false},
		{"long password", "P@ssw0rdP@ssw0rdP@ssw0rdP@ssw0rd12", false},
		{"invalid characters", "P@ssw0rd!", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := utils.ValidatePassword(test.password)
			if err == nil && !test.shouldValidate {
				t.Errorf("expected no error, but got %s", err)
			} else if err != nil && test.shouldValidate {
				t.Error("expected error, but got nil")
			}
		})
	}
}

func TestCreateToken(t *testing.T) {
	tests := []struct {
		name      string
		username  string
		wantErr   bool
		wantToken string
	}{
		{"empty username", "", true, ""},
		{"valid username", "testuser", false, ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			token, err := utils.CreateToken(test.username)
			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, token)
			}
		})
	}
}

func TestVerifyToken(t *testing.T) {
	tests := []struct {
		name        string
		tokenString string
		wantErr     bool
	}{
		{"valid token", createValidToken(t), false},
		{"invalid token", "invalid-token", true},
		{"empty token string", "", true},
		{"expired token", createExpiredToken(t), true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := utils.VerifyToken(test.tokenString)
			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
func createValidToken(t *testing.T) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "testuser",
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(utils.SecretKey)
	require.NoError(t, err)
	return tokenString
}
func createExpiredToken(t *testing.T) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "testuser",
		"exp":      time.Now().Add(-time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(utils.SecretKey)
	require.NoError(t, err)
	return tokenString
}
