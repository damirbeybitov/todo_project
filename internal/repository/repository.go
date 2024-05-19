package repository

import (
	"context"
	"strings"

	"github.com/damirbeybitov/todo_project/internal/models"
	token "github.com/damirbeybitov/todo_project/internal/token"
	pbUser "github.com/damirbeybitov/todo_project/proto/user"
)

type Repository struct {
	MicroServiceClients models.MicroServiceClients
}

func NewRepository(microServiceClients models.MicroServiceClients) *Repository {
	return &Repository{MicroServiceClients: microServiceClients}
}

func (r *Repository) GetUserIdFromRequest(tokenString string) (int64, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	username, err := token.VerifyToken(tokenString)
	if err != nil {
		return 0, err
	}

	response, err := r.MicroServiceClients.UserClient.GetUserIdWithUsername(context.Background(), &pbUser.GetUserIdWithUsernameRequest{
		Username: username,
	})
	if err != nil {
		return 0, err
	}

	return response.Id, nil
}

func (r *Repository) GetUsernameFromRequest(tokenString string) (string, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	username, err := token.VerifyToken(tokenString)
	if err != nil {
		return "", err
	}

	return username, nil
}