package repository

import (
	"database/sql"
	"fmt"

	"github.com/damirbeybitov/todo_project/internal/log"
	token "github.com/damirbeybitov/todo_project/internal/token"
	"golang.org/x/crypto/bcrypt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository (db *sql.DB) *Repository{
	return &Repository{db: db}
}

func (r *Repository) CheckPassword(username string, password string) error {
	var storedPassword string
	err := r.db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedPassword)
	if err != nil {
		log.ErrorLogger.Printf("Failed to retrieve stored password: %v", err)
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		log.ErrorLogger.Printf("Invalid password for user: %s", username)
		return fmt.Errorf("invalid password")
	}

	return nil
}

func (r *Repository) GenerateTokens(username string) (string, string, error) {
	accessToken, err := token.GenerateAccessToken(username)
	if err != nil {
		log.ErrorLogger.Printf("Failed to generate access token: %v", err)
		return "", "", err
	}
	refreshToken, err := token.GenerateRefreshToken(username)
	if err != nil {
		log.ErrorLogger.Printf("Failed to generate refresh token: %v", err)
		return "", "", err
	}

	return accessToken, refreshToken, nil
}