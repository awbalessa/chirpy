package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", fmt.Errorf("Authorization header doesn't exist")
	}

	fields := strings.Fields(auth)
	if len(fields) < 2 || fields[0] != "ApiKey" {
		return "", fmt.Errorf("Authorization header format must be 'ApiKey {KEY}'")
	}

	return fields[1], nil
}

func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", fmt.Errorf("Authorization header doesn't exist")
	}

	fields := strings.Fields(auth)
	if len(fields) < 2 || fields[0] != "Bearer" {
		return "", fmt.Errorf("Authorization header format must be 'Bearer {TOKEN}'")
	}

	return fields[1], nil
}
