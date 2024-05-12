package user

import (
	"context"

	"github.com/damirbeybitov/todo_project/internal/log"
	userPB "github.com/damirbeybitov/todo_project/proto/user"
)

type UserService struct{
	userPB.UnimplementedUserServiceServer
}

func NewUserService() userPB.UserServiceServer {
	return &UserService{}
}

func (s *UserService) RegisterUser(ctx context.Context, req *userPB.RegisterUserRequest) (*userPB.RegisterUserResponse, error) {
	log.InfoLogger.Printf("Registering user with username: %s, email: %s", req.Username, req.Email)

	// Реализация регистрации пользователя

	return &userPB.RegisterUserResponse{
		Id: "123",
	}, nil
}

func (s *UserService) GetUserProfile(ctx context.Context, req *userPB.GetUserProfileRequest) (*userPB.GetUserProfileResponse, error) {
	log.InfoLogger.Printf("Getting user profile for user ID: %s", req.Id)

	// Реализация получения профиля пользователя

	return &userPB.GetUserProfileResponse{
		User: &userPB.User{
			Id:       req.Id,
			Username: "example",
			Email:    "example@example.com",
			// Дополнительные данные профиля пользователя
		},
	}, nil
}