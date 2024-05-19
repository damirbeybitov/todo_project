package repository

import (
	"context"
	"database/sql"
	"encoding/json"
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
	// Insert the task into the database
	result, err := r.db.Exec("INSERT INTO tasks (title, description, status, user_id) VALUES (?, ?, ?, ?)", task.Title, task.Description, task.Status, task.UserId)
	if err != nil {
		log.ErrorLogger.Printf("Failed to create task: %v", err)
		return 0, err
	}

	// Retrieve the last insert ID
	taskID, err := result.LastInsertId()
	if err != nil {
		log.ErrorLogger.Printf("Failed to retrieve last insert ID: %v", err)
		return 0, err
	}

	// Set the task ID
	task.Id = taskID

	// Cache the newly created task in Redis
	taskKey := fmt.Sprintf("task:%d", taskID)
	taskJSON, err := json.Marshal(task)
	if err != nil {
		log.ErrorLogger.Printf("Failed to marshal task for caching: %v", err)
		return taskID, err
	}

	err = r.redis.Set(context.Background(), taskKey, taskJSON, 0).Err()
	if err != nil {
		log.ErrorLogger.Printf("Failed to cache task: %v", err)
		return taskID, err
	}

	// Update the cached list of tasks for the user
	tasksKey := fmt.Sprintf("tasks:user:%d", task.UserId)
	allTasksData, err := r.redis.Get(context.Background(), tasksKey).Result()
	if err == redis.Nil {
		// If tasks not found in cache, initialize the list with the new task
		tasks := []models.Task{task}
		tasksJSON, err := json.Marshal(tasks)
		if err != nil {
			log.ErrorLogger.Printf("Failed to marshal tasks: %v", err)
			return taskID, err
		}
		r.redis.Set(context.Background(), tasksKey, tasksJSON, 0)
		log.InfoLogger.Printf("Tasks list created and cached for user ID: %d", task.UserId)
	} else if err != nil {
		log.ErrorLogger.Printf("Failed to get tasks from cache: %v", err)
		return taskID, err
	} else {
		// If tasks found in cache, update the list with the new task
		var tasks []models.Task
		err = json.Unmarshal([]byte(allTasksData), &tasks)
		if err != nil {
			log.ErrorLogger.Printf("Failed to unmarshal tasks from cache: %v", err)
			return taskID, err
		}
		tasks = append(tasks, task)
		tasksJSON, err := json.Marshal(tasks)
		if err != nil {
			log.ErrorLogger.Printf("Failed to marshal tasks: %v", err)
			return taskID, err
		}
		r.redis.Set(context.Background(), tasksKey, tasksJSON, 0)
		log.InfoLogger.Printf("Tasks list updated and cached for user ID: %d", task.UserId)
	}

	log.InfoLogger.Printf("Task created and cached with ID: %d", taskID)
	return taskID, nil
}

func (r *Repository) GetTaskByID(taskID int64) (models.Task, error) {
	var task models.Task

	taskKey := fmt.Sprintf("task:%d", taskID)
	taskData, err := r.redis.Get(context.Background(), taskKey).Result()
	if err == redis.Nil {
		// If task not found in cache, get it from the database
		log.InfoLogger.Printf("Task not found in cache, fetching from database")
		err = r.db.QueryRow("SELECT id, title, description, status, user_id FROM tasks WHERE id = ?", taskID).
			Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.UserId)
		if err != nil {
			log.ErrorLogger.Printf("Failed to get task from db: %v", err)
			return task, err
		}

		// Cache the task in Redis
		taskJSON, _ := json.Marshal(task)
		r.redis.Set(context.Background(), taskKey, taskJSON, 0)
		log.InfoLogger.Printf("Task cached: %s", taskJSON)
	} else if err != nil {
		log.ErrorLogger.Printf("Failed to get task from cache: %v", err)
		return task, err
	} else {
		log.InfoLogger.Printf("Task found in cache: %s", taskData)
		json.Unmarshal([]byte(taskData), &task)
	}

	return task, nil
}

func (r *Repository) GetTasks(userID int64) ([]models.Task, error) {
	var tasks []models.Task

	tasksKey := fmt.Sprintf("tasks:user:%d", userID)
	allTasksData, err := r.redis.Get(context.Background(), tasksKey).Result()
	if err == redis.Nil {
		// If tasks not found in cache, get them from the database
		log.InfoLogger.Printf("Tasks not found in cache, fetching from database")
		rows, err := r.db.Query("SELECT id, title, description, status, user_id FROM tasks WHERE user_id = ?", userID)
		if err != nil {
			log.ErrorLogger.Printf("Failed to get tasks from db: %v", err)
			return tasks, err
		}
		defer rows.Close()

		for rows.Next() {
			var task models.Task
			if err := rows.Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.UserId); err != nil {
				log.ErrorLogger.Printf("Failed to scan task: %v", err)
				return tasks, err
			}
			tasks = append(tasks, task)
		}
		if err = rows.Err(); err != nil {
			log.ErrorLogger.Printf("Rows error: %v", err)
			return tasks, err
		}

		// Cache the tasks in Redis
		tasksJSON, err := json.Marshal(tasks)
		if err != nil {
			log.ErrorLogger.Printf("Failed to marshal tasks: %v", err)
			return tasks, err
		}

		r.redis.Set(context.Background(), tasksKey, tasksJSON, 0)
		log.InfoLogger.Printf("Tasks cached: %s", tasksJSON)
	} else if err != nil {
		log.ErrorLogger.Printf("Failed to get tasks from cache: %v", err)
		return tasks, err
	} else {
		log.InfoLogger.Printf("Tasks found in cache: %s", allTasksData)
		json.Unmarshal([]byte(allTasksData), &tasks)
	}

	return tasks, nil
}

func (r *Repository) UpdateTask(task models.Task) error {
	_, err := r.db.Exec("UPDATE tasks SET title = ?, description = ?, status = ? WHERE id = ?", task.Title, task.Description, task.Status, task.Id)
	if err != nil {
		log.ErrorLogger.Printf("Failed to update task: %v", err)
		return err
	}

	// Update the task cache in Redis
	taskKey := fmt.Sprintf("task:%d", task.Id)
	taskJSON, err := json.Marshal(task)
	if err != nil {
		log.ErrorLogger.Printf("Failed to marshal task for caching: %v", err)
		return err
	}

	err = r.redis.Set(context.Background(), taskKey, taskJSON, 0).Err()
	if err != nil {
		log.ErrorLogger.Printf("Failed to update task cache: %v", err)
		return err
	}

	// Update the list of tasks for the user in Redis
	tasksKey := fmt.Sprintf("tasks:user:%d", task.UserId)
	allTasksData, err := r.redis.Get(context.Background(), tasksKey).Result()
	if err == redis.Nil {
		// If the list is not in cache, skip updating (as it would be re-cached on next retrieval)
		log.InfoLogger.Printf("User's task list not in cache, skipping update")
	} else if err != nil {
		log.ErrorLogger.Printf("Failed to get user's task list from cache: %v", err)
		return err
	} else {
		// Update the cached list of tasks
		var tasks []models.Task
		err = json.Unmarshal([]byte(allTasksData), &tasks)
		if err != nil {
			log.ErrorLogger.Printf("Failed to unmarshal tasks from cache: %v", err)
			return err
		}
		for i, t := range tasks {
			if t.Id == task.Id {
				tasks[i] = task
				break
			}
		}
		tasksJSON, err := json.Marshal(tasks)
		if err != nil {
			log.ErrorLogger.Printf("Failed to marshal updated tasks: %v", err)
			return err
		}
		r.redis.Set(context.Background(), tasksKey, tasksJSON, 0)
		log.InfoLogger.Printf("User's task list updated in cache")
	}

	log.InfoLogger.Printf("Task updated and cached with ID: %d", task.Id)
	return nil
}

func (r *Repository) DeleteTask(taskID int64) error {
	// Get the task to retrieve userID
	var userID int64
	err := r.db.QueryRow("SELECT user_id FROM tasks WHERE id = ?", taskID).Scan(&userID)
	if err != nil {
		log.ErrorLogger.Printf("Failed to get task user_id: %v", err)
		return err
	}

	// Delete the task from the database
	_, err = r.db.Exec("DELETE FROM tasks WHERE id = ?", taskID)
	if err != nil {
		log.ErrorLogger.Printf("Failed to delete task: %v", err)
		return err
	}

	// Delete the task cache in Redis
	taskKey := fmt.Sprintf("task:%d", taskID)
	err = r.redis.Del(context.Background(), taskKey).Err()
	if err != nil {
		log.ErrorLogger.Printf("Failed to delete task cache: %v", err)
		return err
	}

	// Update the list of tasks for the user in Redis
	tasksKey := fmt.Sprintf("tasks:user:%d", userID)
	allTasksData, err := r.redis.Get(context.Background(), tasksKey).Result()
	if err == redis.Nil {
		// If the list is not in cache, skip updating (as it would be re-cached on next retrieval)
		log.InfoLogger.Printf("User's task list not in cache, skipping update")
	} else if err != nil {
		log.ErrorLogger.Printf("Failed to get user's task list from cache: %v", err)
		return err
	} else {
		// Update the cached list of tasks
		var tasks []models.Task
		err = json.Unmarshal([]byte(allTasksData), &tasks)
		if err != nil {
			log.ErrorLogger.Printf("Failed to unmarshal tasks from cache: %v", err)
			return err
		}
		for i, t := range tasks {
			if t.Id == taskID {
				tasks = append(tasks[:i], tasks[i+1:]...)
				break
			}
		}
		tasksJSON, err := json.Marshal(tasks)
		if err != nil {
			log.ErrorLogger.Printf("Failed to marshal updated tasks: %v", err)
			return err
		}
		r.redis.Set(context.Background(), tasksKey, tasksJSON, 0)
		log.InfoLogger.Printf("User's task list updated in cache after deletion")
	}

	log.InfoLogger.Printf("Task deleted and cache updated with ID: %d", taskID)
	return nil
}

func (r *Repository) GetUserIdWithUsername(username string) (int64, error) {
	var id int64
	err := r.db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&id)
	if err != nil {
		log.ErrorLogger.Printf("Failed to get user ID: %v", err)
		return 0, err
	}

	return id, nil
}
