package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/damirbeybitov/todo_project/internal/log"
	"golang.org/x/crypto/bcrypt"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) CheckUserInDB(tx *sql.Tx, username string, email string) error {
	// Check if the user already exists
	var count int
	err := tx.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM users WHERE username = ? OR email = ?", username, email).Scan(&count)
	if err != nil {
		tx.Rollback()
		log.ErrorLogger.Printf("Failed to check user existence: %v", err)
		return err
	}

	if count > 0 {
		tx.Rollback()
		log.ErrorLogger.Printf("Username or email already exists")
		return fmt.Errorf("username or email already exists")
	}

	return nil
}

func (r *Repository) AddUserToDB(tx *sql.Tx, username string, email string, password string) (int64, error) {
	result, err := tx.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, password)
	if err != nil {
		tx.Rollback()
		log.ErrorLogger.Printf("Failed to insert user: %v", err)
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.ErrorLogger.Printf("Failed to get last insert ID: %v", err)
		return -1, err
	}

	return id, nil
}

func (r *Repository) CheckPassword(username string, password string) error {
	var storedPassword string
	err := r.DB.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedPassword)
	if err != nil {
		log.ErrorLogger.Printf("Failed to retrieve stored password: %v", err)
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		log.ErrorLogger.Printf("Invalid password for user %s: %v", username, err)
		return err
	}

	return nil
}

func (r *Repository) DeleteUserFromDB(tx *sql.Tx, username string) error {
	result, err := r.DB.Exec("DELETE FROM users WHERE username = ?", username)
	if err != nil {
		log.ErrorLogger.Printf("Failed to delete user: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.ErrorLogger.Printf("Failed to get rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		log.ErrorLogger.Printf("User not found")
		return fmt.Errorf("user not found")
	}

	return nil
}
