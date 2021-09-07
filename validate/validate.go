package validate

import (
	"github.com/Bruary/twitter-clone/service/models"
	"github.com/dgrijalva/jwt-go"
)

const PasswordMinLength = 8

var jwtKey = []byte("White_Yasmin")

func IsStringEmpty(text string) bool {
	return text == ""
}

func IsNumberNegative(number int) bool {
	return number < 0
}

func IsAge12AndAbove(number int) bool {
	return number >= 12
}

func IsPasswordLengthCorrect(password string) bool {
	return len(password) >= PasswordMinLength
}

func IsTokenValid(tokenString string) bool {
	claims := &models.Claims{}
	token, _ := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token.Valid
}

func GetJWTclaims(tokenString string) *models.Claims {

	claims := &models.Claims{}
	jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	return claims
}
