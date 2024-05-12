package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/damirbeybitov/todo_project/internal/log"
	userPB "github.com/damirbeybitov/todo_project/proto/user"
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

	// Insert the new user
	result, err := tx.ExecContext(ctx, "INSERT INTO users (username, email, password) VALUES (?, ?, ?)", req.Username, req.Email, req.Password)
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