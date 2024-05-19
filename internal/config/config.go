package config

import (
	"encoding/json"
	"os"

	"github.com/damirbeybitov/todo_project/internal/models"
)

func NewConfig(fileName string) (*models.Config, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var Config models.Config
	err = json.NewDecoder(file).Decode(&Config)
	if err != nil {
		return nil, err
	}

	return &Config, nil
}