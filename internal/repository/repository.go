package repository

import "github.com/damirbeybitov/todo_project/internal/models"

type Repository struct {
	MicroServiceClients models.MicroServiceClients
}

func NewRepository(microServiceClients models.MicroServiceClients) *Repository {
	return &Repository{MicroServiceClients: microServiceClients}
}