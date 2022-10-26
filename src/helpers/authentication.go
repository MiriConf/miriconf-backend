package helpers

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hashedPassword), err
}

func CheckPassword(hashedPassword, storedPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(storedPass))
	return err == nil
}
