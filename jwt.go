package main

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("secret-key") //секретный ключ для подписи токенов

// Структура для хранения данных о пользователе
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func generateToken(username string) (string, error) {
	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func validateToken(token string) bool {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return false
	}
	expiresString := parsed.Claims.(jwt.MapClaims)["Expires"].(string)
	expires, err := strconv.ParseInt(expiresString, 10, 64)
	if err != nil {
		return false
	}
	if time.Now().Unix() > expires {
		return false
	}
	return err == nil
}
