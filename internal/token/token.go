package auth

import (
	"errors"
	"time"

	"github.com/damirbeybitov/todo_project/internal/log"
	"github.com/dgrijalva/jwt-go"
)

const (
	accessTokenTime 	= time.Minute * 15
	refreshTokenTime 	= time.Hour * 24
	signingKey			= "yIAYiuIoibngJG78G785F76"
)

func GenerateRefreshToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(refreshTokenTime).Unix(),
		Subject:   username,
	})

	return token.SignedString([]byte(signingKey))
}

func GenerateAccessToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(accessTokenTime).Unix(),
		Subject:   username,
	})

	return token.SignedString([]byte(signingKey))
}

func VerifyToken(token string) (string, error){
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})
	if err != nil {
		log.ErrorLogger.Printf("Error parsing JWT token: %v", err)
		return "", err
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok || !t.Valid {
		log.ErrorLogger.Printf("Invalid JWT token")
		return "", errors.New("invalid JWT token")
	}

	return claims["sub"].(string), nil
}