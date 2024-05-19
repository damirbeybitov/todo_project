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
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
	UserId      int64  `json:"user_id"`
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
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type DeleteUserRequest struct {
	Password string `json:"password"`
}

type DeleteUserResponse struct {
	Message string `json:"message"`
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	UserId      int64  `json:"user_id"`
}

type CreateTaskResponse struct {
	Id int64 `json:"id"`
}

type GetTasksResponse struct {
	Tasks []Task `json:"tasks"`
}

type UpdateTaskRequest struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
	UserId      int64  `json:"user_id"`
}

type UpdateTaskResponse struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
	UserId      int64  `json:"user_id"`
}

type DeleteTaskRequest struct {
	Id int64 `json:"id"`
}

type DeleteTaskResponse struct {
	Message string `json:"message"`
}
