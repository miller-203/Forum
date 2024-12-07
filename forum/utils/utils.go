package forum

import (
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
func VerifyPassword(hashedPassword, plainPassword string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {

		return false
	}

	return true
}

func IsValidUsername(username string) bool {

	// If the username is empty after trimming spaces, it's invalid
	if username == "" {
		return false
	}

	// Check each character in the username
	for _, char := range username {

		if !(unicode.IsLetter(char) || unicode.IsDigit(char)) {
			return false
		}
	}

	// If all characters are valid, return true
	return true
}
