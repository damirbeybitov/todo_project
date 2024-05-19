package repository

import (
	"database/sql"
	"fmt"

	"github.com/damirbeybitov/todo_project/internal/log"
	"github.com/damirbeybitov/todo_project/internal/models"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewRepository(db *sql.DB, redis *redis.Client) *Repository {
	return &Repository{db: db, redis: redis}
}

func (r *Repository) CreateTask(task models.Task) (int64, error) {
	result, err := r.db.Exec("INSERT INTO tasks (title, description, status, user_id) VALUES (?, ?, ?, ?)", task.Title, task.Description, task.Status, task.UserID)
	if err != nil {
		log.ErrorLogger.Printf("Failed to create task: %v", err)
		return 0, err
	}

	taskID, err := result.LastInsertId()
	if err != nil {
		log.ErrorLogger.Printf("Failed to retrieve last insert ID: %v", err)
		return 0, err
	}

	return taskID, nil
}

func (r *Repository) GetTaskByID(taskID int64) (models.Task, error) {
	task := models.Task{}
	err := r.db.QueryRow("SELECT id, title, description, status, user_id FROM tasks WHERE id = ?", taskID).Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return task, fmt.Errorf("task not found")
		}
		log.ErrorLogger.Printf("Failed to get task by ID: %v", err)
		return task, err
	}

	return task, nil
}

func (r *Repository) UpdateTask(task models.Task) error {
	_, err := r.db.Exec("UPDATE tasks SET title = ?, description = ?, status = ? WHERE id = ?", task.Title, task.Description, task.Status, task.ID)
	if err != nil {
		log.ErrorLogger.Printf("Failed to update task: %v", err)
		return err
	}

	return nil
}

func (r *Repository) DeleteTask(taskID int64) error {
	_, err := r.db.Exec("DELETE FROM tasks WHERE id = ?", taskID)
	if err != nil {
		log.ErrorLogger.Printf("Failed to delete task: %v", err)
		return err
	}

	return nil
}
