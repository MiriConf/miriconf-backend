package helpers

import (
	"os"
	"strings"

	jwt "github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Token struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func ValidateToken(headerToken string) (*jwt.Token, error) {
	jwtKey := []byte(os.Getenv("JWT_KEY"))
	tokenString := strings.Split(headerToken, " ")[1]
	token, err := jwt.ParseWithClaims(tokenString, &Token{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token, err
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hashedPassword), err
}

func CheckPassword(hashedPassword, storedPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(storedPass))
	return err == nil
}
