package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/damirbeybitov/todo_project/internal/log"
	userPB "github.com/damirbeybitov/todo_project/proto/user"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct{
	db *sql.DB
	userPB.UnimplementedUserServiceServer
}

func NewUserService(db *sql.DB) userPB.UserServiceServer {
	return &UserService{db: db}
}

func (s *UserService) RegisterUser(ctx context.Context, req *userPB.RegisterUserRequest) (*userPB.RegisterUserResponse, error) {
	log.InfoLogger.Printf("Registering user with username: %s, email: %s", req.Username, req.Email)

	// Реализация регистрации пользователя
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.ErrorLogger.Printf("Failed to start transaction: %v", err)
		return nil, err
	}

	// Check if the user already exists
	var count int
	err = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE username = ? OR email = ?", req.Username, req.Email).Scan(&count)
	if err != nil {
		tx.Rollback()
		log.ErrorLogger.Printf("Failed to check user existence: %v", err)
		return nil, err
	}

	if count > 0 {
		tx.Rollback()
		log.ErrorLogger.Printf("Username or email already exists")
		return nil, fmt.Errorf("username or email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		log.ErrorLogger.Printf("Failed to hash password: %v", err)
		return nil, err
	}

	hashedPasswordStr := string(hashedPassword)
	
	// Insert the new user
	result, err := tx.ExecContext(ctx, "INSERT INTO users (username, email, password) VALUES (?, ?, ?)", req.Username, req.Email, hashedPasswordStr)
	if err != nil {
		tx.Rollback()
		log.ErrorLogger.Printf("Failed to insert user: %v", err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.ErrorLogger.Printf("Failed to get last insert ID: %v", err)
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.ErrorLogger.Printf("Failed to commit transaction: %v", err)
		return nil, err
	}

	log.InfoLogger.Printf("User registered successfully with ID: %d", id)

	return &userPB.RegisterUserResponse{
		Id: fmt.Sprintf("%d", id),
	}, nil
}

func (s *UserService) GetUserProfile(ctx context.Context, req *userPB.GetUserProfileRequest) (*userPB.GetUserProfileResponse, error) {
	log.InfoLogger.Printf("Getting user profile for user ID: %s", req.Id)

	// Реализация получения профиля пользователя
	var username, email string
	err := s.db.QueryRowContext(ctx, "SELECT username, email FROM users WHERE id = ?", req.Id).Scan(&username, &email)
	if err != nil {
		log.ErrorLogger.Printf("Failed to get user profile: %v", err)
		return nil, err
	}
	log.InfoLogger.Printf("User profile retrieved successfully for user ID: %s", req.Id)

	return &userPB.GetUserProfileResponse{
		User: &userPB.User{
			Id:       req.Id,
			Username: username,
			Email:    email,
		},
	}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *userPB.DeleteUserRequest) (*userPB.DeleteUserResponse, error) {
	log.InfoLogger.Printf("Deleting user with ID: %s", req.Username)

	// Check if the provided password matches the username
	var storedPassword string
	err := s.db.QueryRowContext(ctx, "SELECT password FROM users WHERE username = ?", req.Username).Scan(&storedPassword)
	if err != nil {
		log.ErrorLogger.Printf("Failed to retrieve stored password: %v", err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(req.Password))
	if err != nil {
		log.ErrorLogger.Printf("Invalid password for user: %s", req.Username)
		return nil, fmt.Errorf("invalid password")
	}

	// Реализация удаления пользователя
	result, err := s.db.ExecContext(ctx, "DELETE FROM users WHERE username = ?", req.Username)
	if err != nil {
		log.ErrorLogger.Printf("Failed to delete user: %v", err)
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.ErrorLogger.Printf("Failed to get rows affected: %v", err)
		return nil, err
	}

	if rowsAffected == 0 {
		log.ErrorLogger.Printf("User not found")
		return nil, fmt.Errorf("user not found")
	}

	log.InfoLogger.Printf("User deleted successfully with username: %s", req.Username)

	message := fmt.Sprintf("User deleted successfully with username: %s", req.Username)
	return &userPB.DeleteUserResponse{
		Message: message,
	}, nil
}