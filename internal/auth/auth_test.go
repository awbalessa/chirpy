package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	tests := []struct {
		name        string
		userID      uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
	}{
		{
			name:        "Valid token with 1 hour expiry",
			userID:      uuid.New(),
			tokenSecret: "test-secret-key",
			expiresIn:   time.Hour,
		},
		{
			name:        "Valid token with 24 hour expiry",
			userID:      uuid.New(),
			tokenSecret: "another-test-secret",
			expiresIn:   time.Hour * 24,
		},
		{
			name:        "Valid token with 1 minute expiry",
			userID:      uuid.New(),
			tokenSecret: "one-more-test-secret",
			expiresIn:   time.Minute,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			token, err := MakeJWT(tc.userID, tc.tokenSecret, tc.expiresIn)
			if err != nil {
				t.Fatalf("Error creating token: %v", err)
			}

			gotUserID, err := ValidateJWT(token, tc.tokenSecret)
			if err != nil {
				t.Fatalf("Error validating JWT: %v", err)
			}

			if gotUserID != tc.userID {
				t.Fatalf("Expected user ID %v, got %v", tc.userID, gotUserID)
			}
		})
	}
}

func TestExpiredToken(t *testing.T) {
	test := struct {
		name        string
		userID      uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
	}{
		name:        "Expired token",
		userID:      uuid.New(),
		tokenSecret: "my-top-secret",
		expiresIn:   -1 * time.Second,
	}

	t.Run(test.name, func(t *testing.T) {
		token, err := MakeJWT(test.userID, test.tokenSecret, test.expiresIn)
		if err != nil {
			t.Fatalf("Error making JWT: %v", err)
		}

		_, err = ValidateJWT(token, test.tokenSecret)
		if err == nil {
			t.Fatal("Expected error for expired token, got nil")
		}
	})
}

func TestWrongToken(t *testing.T) {
	test := struct {
		name        string
		userID      uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
	}{
		name:        "Expired token",
		userID:      uuid.New(),
		tokenSecret: "correct-secret",
		expiresIn:   time.Hour,
	}

	t.Run(test.name, func(t *testing.T) {
		token, err := MakeJWT(test.userID, test.tokenSecret, test.expiresIn)
		if err != nil {
			t.Fatalf("Error making JWT: %v", err)
		}

		_, err = ValidateJWT(token, "wrong-secret")
		if err == nil {
			t.Fatal("Expected wrong token error, got nil")
		}
	})
}
