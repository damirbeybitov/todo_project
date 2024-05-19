package models

import (
	pbAuth "github.com/damirbeybitov/todo_project/proto/auth"
	pbTask "github.com/damirbeybitov/todo_project/proto/task"
	pbUser "github.com/damirbeybitov/todo_project/proto/user"
)

type Config struct {
	SqlConnection string `json:"sqlConnection"`
}

type Task struct {
	ID          int64
	Title       string
	Description string
	Status      bool
	UserID      int64
}

type MicroServiceClients struct {
	UserClient pbUser.UserServiceClient
	AuthClient pbAuth.AuthServiceClient
	TaskClient pbTask.TaskServiceClient
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Id int64 `json:"id"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type GetUserProfileRequest struct {
	Id int64 `json:"id"`
}

type GetUserProfileResponse struct {
	Id int64 `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type DeleteUserRequest struct {
	Password string `json:"password"`
}

type DeleteUserResponse struct {
	Message string `json:"message"`
}