package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	}
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := newToken.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("Error signing token: %v", err)
	}

	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	// Define a claims struct to parse into
	claims := &jwt.RegisteredClaims{}

	// Parse the token
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			// Validate the signing method is what we expect
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Return the secret key used to sign
			return []byte(tokenSecret), nil
		},
	)

	// Check for parsing errors
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse token: %v", err)
	}

	// Check if token is valid
	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	// Extract the subject (user ID)
	subject := claims.Subject

	// Convert the string ID to UUID
	userID, err := uuid.Parse(subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID in token: %v", err)
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", fmt.Errorf("Authorization header doesn't exist")
	}

	fields := strings.Fields(auth)
	if len(fields) < 2 || fields[0] != "Bearer" {
		return "", fmt.Errorf("Authorization header format must be 'Bearer {token}'")
	}

	return fields[1], nil
}
