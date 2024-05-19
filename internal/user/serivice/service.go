package user

import (
	"context"
	"fmt"

	"github.com/damirbeybitov/todo_project/internal/log"
	"github.com/damirbeybitov/todo_project/internal/user/repository"
	userPB "github.com/damirbeybitov/todo_project/proto/user"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct{
	repo *repository.Repository
	userPB.UnimplementedUserServiceServer
}

func NewUserService(repo *repository.Repository) userPB.UserServiceServer {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(ctx context.Context, req *userPB.RegisterUserRequest) (*userPB.RegisterUserResponse, error) {
	log.InfoLogger.Printf("Registering user with username: %s, email: %s", req.Username, req.Email)

	// Реализация регистрации пользователя
	tx, err := s.repo.DB.BeginTx(ctx, nil)
	if err != nil {
		log.ErrorLogger.Printf("Failed to start transaction: %v", err)
		return nil, err
	}

	// Check if the user already exists
	if err := s.repo.CheckUserInDB(tx, req.Username, req.Email); err != nil {
		return nil, err
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
	id, err := s.repo.AddUserToDB(tx, req.Username, req.Email, hashedPasswordStr)
	if err != nil {
		return nil, err
	
	}

	err = tx.Commit()
	if err != nil {
		log.ErrorLogger.Printf("Failed to commit transaction: %v", err)
		return nil, err
	}

	log.InfoLogger.Printf("User %s registered successfully with ID: %d", req.Username, id)

	return &userPB.RegisterUserResponse{
		Id: id,
	}, nil
}

func (s *UserService) GetUserProfile(ctx context.Context, req *userPB.GetUserProfileRequest) (*userPB.GetUserProfileResponse, error) {
	log.InfoLogger.Printf("Getting user profile for user ID: %d", req.Id)

	// Реализация получения профиля пользователя
	var username, email string
	err := s.repo.DB.QueryRowContext(ctx, "SELECT username, email FROM users WHERE id = ?", req.Id).Scan(&username, &email)
	if err != nil {
		log.ErrorLogger.Printf("Failed to get user profile: %v", err)
		return nil, err
	}

	log.InfoLogger.Printf("User profile retrieved successfully for user ID: %d", req.Id)

	return &userPB.GetUserProfileResponse{
		User: &userPB.User{
			Id:       req.Id,
			Username: username,
			Email:    email,
		},
	}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *userPB.DeleteUserRequest) (*userPB.DeleteUserResponse, error) {
	log.InfoLogger.Printf("Deleting user with username: %s", req.Username)

	// Check if the provided password matches the username
	if err := s.repo.CheckPassword(req.Username, req.Password); err != nil {
		return nil, err
	}

	// Реализация удаления пользователя
	if err := s.repo.DeleteuserFromDB(req.Username); err != nil {
		return nil, err
	}

	log.InfoLogger.Printf("User deleted successfully with username: %s", req.Username)

	message := fmt.Sprintf("User deleted successfully with username: %s", req.Username)
	return &userPB.DeleteUserResponse{
		Message: message,
	}, nil
}

func (s *UserService) GetUserIdWithUsername(ctx context.Context, req *userPB.GetUserIdWithUsernameRequest) (*userPB.GetUserIdWithUsernameResponse, error) {
	log.InfoLogger.Printf("Getting user ID for username: %s", req.Username)

	// Реализация получения ID пользователя по имени
	var id int64
	err := s.repo.DB.QueryRowContext(ctx, "SELECT id FROM users WHERE username = ?", req.Username).Scan(&id)
	if err != nil {
		log.ErrorLogger.Printf("Failed to get user ID: %v", err)
		return nil, err
	}

	log.InfoLogger.Printf("User ID retrieved successfully for username: %s", req.Username)

	return &userPB.GetUserIdWithUsernameResponse{
		Id: id,
	}, nil
}