package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pwd string) (string, error) {
	bytedPwd := []byte(pwd)
	hashedBytes, err := bcrypt.GenerateFromPassword(bytedPwd, bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Error generating hash: %v", err)
	}

	return string(hashedBytes), nil
}

func CheckPasswordHash(hash, pwd string) error {
	bytedPwd := []byte(pwd)
	bytedHash := []byte(hash)
	return bcrypt.CompareHashAndPassword(bytedHash, bytedPwd)
}
