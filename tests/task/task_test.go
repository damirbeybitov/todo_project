package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/damirbeybitov/todo_project/internal/models"
	"github.com/damirbeybitov/todo_project/internal/task/repository"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	// Setup the database connection
	dsn := "root:root@tcp(localhost:3306)/to_do"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Setup the Redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()

	// Ensure Redis is empty before starting the test
	err = rdb.FlushAll(context.Background()).Err()
	if err != nil {
		t.Fatalf("Failed to flush Redis: %v", err)
	}

	// Create the repository
	repo := repository.NewRepository(db, rdb)

	// Create a task
	task := models.Task{
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      true,
		UserId:      1,
	}

	// Call the function to test
	taskID, err := repo.CreateTask(task)
	assert.NoError(t, err, "Expected no error from CreateTask")
	assert.NotZero(t, taskID, "Expected task ID to be non-zero")

	// Verify the task was inserted into the database
	var retrievedTask models.Task
	err = db.QueryRow("SELECT id, title, description, status, user_id FROM tasks WHERE id = ?", taskID).
		Scan(&retrievedTask.Id, &retrievedTask.Title, &retrievedTask.Description, &retrievedTask.Status, &retrievedTask.UserId)
	assert.NoError(t, err, "Expected no error when querying the task from the database")
	assert.Equal(t, task.Title, retrievedTask.Title, "Expected task title to match")
	assert.Equal(t, task.Description, retrievedTask.Description, "Expected task description to match")
	assert.Equal(t, task.Status, retrievedTask.Status, "Expected task status to match")
	assert.Equal(t, task.UserId, retrievedTask.UserId, "Expected task user ID to match")

	// Verify the task was cached in Redis
	taskKey := fmt.Sprintf("task:%d", taskID)
	taskJSON, err := rdb.Get(context.Background(), taskKey).Result()
	assert.NoError(t, err, "Expected no error when getting the task from Redis")

	var cachedTask models.Task
	err = json.Unmarshal([]byte(taskJSON), &cachedTask)
	assert.NoError(t, err, "Expected no error when unmarshaling the cached task")
	assert.Equal(t, task.Title, cachedTask.Title, "Expected cached task title to match")
	assert.Equal(t, task.Description, cachedTask.Description, "Expected cached task description to match")
	assert.Equal(t, task.Status, cachedTask.Status, "Expected cached task status to match")
	assert.Equal(t, task.UserId, cachedTask.UserId, "Expected cached task user ID to match")

	// Verify the task list was updated in Redis
	tasksKey := fmt.Sprintf("tasks:user:%d", task.UserId)
	tasksJSON, err := rdb.Get(context.Background(), tasksKey).Result()
	assert.NoError(t, err, "Expected no error when getting the task list from Redis")

	var tasks []models.Task
	err = json.Unmarshal([]byte(tasksJSON), &tasks)
	assert.NoError(t, err, "Expected no error when unmarshaling the task list")
	assert.Len(t, tasks, 1, "Expected the task list to contain one task")
	assert.Equal(t, taskID, tasks[0].Id, "Expected the task list to contain the correct task ID")
}

func TestGetTaskByID(t *testing.T) {
	// Setup the database connection
	dsn := "root:root@tcp(localhost:3306)/to_do"
	db, err := sql.Open("mysql", dsn)
	assert.NoError(t, err, "Failed to connect to the database")
	defer db.Close()

	// Setup the Redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()

	// Ensure Redis is empty before starting the test
	err = rdb.FlushAll(context.Background()).Err()
	assert.NoError(t, err, "Failed to flush Redis")

	// Create the repository
	repo := repository.NewRepository(db, rdb)

	// Create a task to retrieve
	task := models.Task{
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      true,
		UserId:      1,
	}

	// Insert the task into the database directly for testing
	result, err := db.Exec("INSERT INTO tasks (title, description, status, user_id) VALUES (?, ?, ?, ?)", task.Title, task.Description, task.Status, task.UserId)
	assert.NoError(t, err, "Failed to insert task into the database")

	taskID, err := result.LastInsertId()
	assert.NoError(t, err, "Failed to get the last insert ID")
	task.Id = taskID

	// Test retrieval from the database when cache is empty
	retrievedTask, err := repo.GetTaskByID(taskID)
	assert.NoError(t, err, "Expected no error from GetTaskByID")
	assert.Equal(t, task, retrievedTask, "Expected retrieved task to match the inserted task")

	// Verify the task was cached in Redis
	taskKey := fmt.Sprintf("task:%d", taskID)
	taskJSON, err := rdb.Get(context.Background(), taskKey).Result()
	assert.NoError(t, err, "Expected no error when getting the task from Redis")

	var cachedTask models.Task
	err = json.Unmarshal([]byte(taskJSON), &cachedTask)
	assert.NoError(t, err, "Expected no error when unmarshaling the cached task")
	assert.Equal(t, task, cachedTask, "Expected cached task to match the inserted task")

	// Test retrieval from the cache
	retrievedTask, err = repo.GetTaskByID(taskID)
	assert.NoError(t, err, "Expected no error from GetTaskByID on cache hit")
	assert.Equal(t, task, retrievedTask, "Expected retrieved task to match the cached task")
}

func TestGetTasks(t *testing.T) {
	// Setup the database connection
	dsn := "root:root@tcp(localhost:3306)/to_do"
	db, err := sql.Open("mysql", dsn)
	assert.NoError(t, err, "Failed to connect to the database")
	defer db.Close()

	// Setup the Redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()

	// Ensure Redis is empty before starting the test
	err = rdb.FlushAll(context.Background()).Err()
	assert.NoError(t, err, "Failed to flush Redis")

	// Create the repository
	repo := repository.NewRepository(db, rdb)

	// Create tasks to retrieve
	tasks := []models.Task{
		{Title: "Task 1", Description: "Description 1", Status: false, UserId: 1},
		{Title: "Task 2", Description: "Description 2", Status: true, UserId: 1},
	}

	// Insert the tasks into the database directly for testing
	for _, task := range tasks {
		result, err := db.Exec("INSERT INTO tasks (title, description, status, user_id) VALUES (?, ?, ?, ?)", task.Title, task.Description, task.Status, task.UserId)
		assert.NoError(t, err, "Failed to insert task into the database")
		taskID, err := result.LastInsertId()
		assert.NoError(t, err, "Failed to get the last insert ID")
		task.Id = taskID
	}

	// Test retrieval from the database when cache is empty
	retrievedTasks, err := repo.GetTasks(1)
	assert.NoError(t, err, "Expected no error from GetTasks")
	assert.Len(t, retrievedTasks, len(tasks), "Expected number of retrieved tasks to match the inserted tasks")
	for i, task := range tasks {
		assert.Equal(t, task.Title, retrievedTasks[i].Title, "Expected task title to match")
		assert.Equal(t, task.Description, retrievedTasks[i].Description, "Expected task description to match")
		assert.Equal(t, task.Status, retrievedTasks[i].Status, "Expected task status to match")
		assert.Equal(t, task.UserId, retrievedTasks[i].UserId, "Expected task user ID to match")
	}

	// Verify the tasks were cached in Redis
	tasksKey := fmt.Sprintf("tasks:user:%d", 1)
	tasksJSON, err := rdb.Get(context.Background(), tasksKey).Result()
	assert.NoError(t, err, "Expected no error when getting the tasks from Redis")

	var cachedTasks []models.Task
	err = json.Unmarshal([]byte(tasksJSON), &cachedTasks)
	assert.NoError(t, err, "Expected no error when unmarshaling the cached tasks")
	assert.Len(t, cachedTasks, len(tasks), "Expected number of cached tasks to match the inserted tasks")

	// Test retrieval from the cache
	retrievedTasks, err = repo.GetTasks(1)
	assert.NoError(t, err, "Expected no error from GetTasks on cache hit")
	assert.Len(t, retrievedTasks, len(tasks), "Expected number of retrieved tasks to match the cached tasks")
	for i, task := range tasks {
		assert.Equal(t, task.Title, retrievedTasks[i].Title, "Expected task title to match")
		assert.Equal(t, task.Description, retrievedTasks[i].Description, "Expected task description to match")
		assert.Equal(t, task.Status, retrievedTasks[i].Status, "Expected task status to match")
		assert.Equal(t, task.UserId, retrievedTasks[i].UserId, "Expected task user ID to match")
	}
}



func TestUpdateTask(t *testing.T) {
	// Setup the database connection
	dsn := "root:root@tcp(localhost:3306)/to_do"
	db, err := sql.Open("mysql", dsn)
	assert.NoError(t, err, "Failed to connect to the database")
	defer db.Close()

	// Setup the Redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()

	// Ensure Redis is empty before starting the test
	err = rdb.FlushAll(context.Background()).Err()
	assert.NoError(t, err, "Failed to flush Redis")

	// Create the repository
	repo := repository.NewRepository(db, rdb)

	// Create a task to update
	task := models.Task{
		Title:       "Original Task",
		Description: "This is the original task",
		Status:      false,
		UserId:      1,
	}

	// Insert the task into the database directly for testing
	result, err := db.Exec("INSERT INTO tasks (title, description, status, user_id) VALUES (?, ?, ?, ?)", task.Title, task.Description, task.Status, task.UserId)
	assert.NoError(t, err, "Failed to insert task into the database")
	taskID, err := result.LastInsertId()
	assert.NoError(t, err, "Failed to get the last insert ID")
	task.Id = taskID

	// Cache the task in Redis
	taskKey := fmt.Sprintf("task:%d", taskID)
	taskJSON, err := json.Marshal(task)
	assert.NoError(t, err, "Failed to marshal task for caching")
	err = rdb.Set(context.Background(), taskKey, taskJSON, 0).Err()
	assert.NoError(t, err, "Failed to set task in Redis")

	// Create a list of tasks and cache it
	tasks := []models.Task{task}
	tasksKey := fmt.Sprintf("tasks:user:%d", task.UserId)
	tasksJSON, err := json.Marshal(tasks)
	assert.NoError(t, err, "Failed to marshal tasks for caching")
	err = rdb.Set(context.Background(), tasksKey, tasksJSON, 0).Err()
	assert.NoError(t, err, "Failed to set tasks in Redis")

	// Update the task
	updatedTask := models.Task{
		Id:          taskID,
		Title:       "Updated Task",
		Description: "This is the updated task",
		Status:      true,
		UserId:      1,
	}

	// Call the function to test
	err = repo.UpdateTask(updatedTask)
	assert.NoError(t, err, "Expected no error from UpdateTask")

	// Verify the task was updated in the database
	var retrievedTask models.Task
	err = db.QueryRow("SELECT id, title, description, status, user_id FROM tasks WHERE id = ?", taskID).
		Scan(&retrievedTask.Id, &retrievedTask.Title, &retrievedTask.Description, &retrievedTask.Status, &retrievedTask.UserId)
	assert.NoError(t, err, "Expected no error when querying the task from the database")
	assert.Equal(t, updatedTask.Title, retrievedTask.Title, "Expected task title to be updated")
	assert.Equal(t, updatedTask.Description, retrievedTask.Description, "Expected task description to be updated")
	assert.Equal(t, updatedTask.Status, retrievedTask.Status, "Expected task status to be updated")
	assert.Equal(t, updatedTask.UserId, retrievedTask.UserId, "Expected task user ID to remain the same")

	// Verify the task was updated in Redis
	taskJSON, err = rdb.Get(context.Background(), taskKey).Bytes()
	assert.NoError(t, err, "Expected no error when getting the task from Redis")

	var cachedTask models.Task
	err = json.Unmarshal([]byte(taskJSON), &cachedTask)
	assert.NoError(t, err, "Expected no error when unmarshaling the cached task")
	assert.Equal(t, updatedTask.Title, cachedTask.Title, "Expected cached task title to be updated")
	assert.Equal(t, updatedTask.Description, cachedTask.Description, "Expected cached task description to be updated")
	assert.Equal(t, updatedTask.Status, cachedTask.Status, "Expected cached task status to be updated")
	assert.Equal(t, updatedTask.UserId, cachedTask.UserId, "Expected cached task user ID to remain the same")

	// Verify the task list was updated in Redis
	tasksJSON, err = rdb.Get(context.Background(), tasksKey).Bytes()
	assert.NoError(t, err, "Expected no error when getting the tasks list from Redis")

	var cachedTasks []models.Task
	err = json.Unmarshal([]byte(tasksJSON), &cachedTasks)
	assert.NoError(t, err, "Expected no error when unmarshaling the cached tasks")
	assert.Len(t, cachedTasks, 1, "Expected the tasks list to contain one task")
	assert.Equal(t, updatedTask.Id, cachedTasks[0].Id, "Expected the task ID to match")
	assert.Equal(t, updatedTask.Title, cachedTasks[0].Title, "Expected the cached task title to be updated")
	assert.Equal(t, updatedTask.Description, cachedTasks[0].Description, "Expected the cached task description to be updated")
	assert.Equal(t, updatedTask.Status, cachedTasks[0].Status, "Expected the cached task status to be updated")
	assert.Equal(t, updatedTask.UserId, cachedTasks[0].UserId, "Expected the cached task user ID to remain the same")
}

func TestDeleteTask(t *testing.T) {
	// Setup the database connection
	dsn := "root:root@tcp(localhost:3306)/to_do"
	db, err := sql.Open("mysql", dsn)
	assert.NoError(t, err, "Failed to connect to the database")
	defer db.Close()

	// Setup the Redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()

	// Ensure Redis is empty before starting the test
	err = rdb.FlushAll(context.Background()).Err()
	assert.NoError(t, err, "Failed to flush Redis")

	// Create the repository
	repo := repository.NewRepository(db, rdb)

	// Create a task to delete
	task := models.Task{
		Title:       "Task to Delete",
		Description: "This task will be deleted",
		Status:      true,
		UserId:      1,
	}

	// Insert the task into the database directly for testing
	result, err := db.Exec("INSERT INTO tasks (title, description, status, user_id) VALUES (?, ?, ?, ?)", task.Title, task.Description, task.Status, task.UserId)
	assert.NoError(t, err, "Failed to insert task into the database")
	taskID, err := result.LastInsertId()
	assert.NoError(t, err, "Failed to get the last insert ID")
	task.Id = taskID

	// Cache the task in Redis
	taskKey := fmt.Sprintf("task:%d", taskID)
	taskJSON, err := json.Marshal(task)
	assert.NoError(t, err, "Failed to marshal task for caching")
	err = rdb.Set(context.Background(), taskKey, taskJSON, 0).Err()
	assert.NoError(t, err, "Failed to set task in Redis")

	// Create a list of tasks and cache it
	tasks := []models.Task{task}
	tasksKey := fmt.Sprintf("tasks:user:%d", task.UserId)
	tasksJSON, err := json.Marshal(tasks)
	assert.NoError(t, err, "Failed to marshal tasks for caching")
	err = rdb.Set(context.Background(), tasksKey, tasksJSON, 0).Err()
	assert.NoError(t, err, "Failed to set tasks in Redis")

	// Call the function to test
	err = repo.DeleteTask(taskID)
	assert.NoError(t, err, "Expected no error from DeleteTask")

	// Verify the task was deleted from the database
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM tasks WHERE id = ?", taskID).Scan(&count)
	assert.NoError(t, err, "Expected no error when querying the task from the database")
	assert.Equal(t, 0, count, "Expected the task to be deleted from the database")

	// Verify the task was deleted from Redis
	taskJSON, err = rdb.Get(context.Background(), taskKey).Bytes()
	assert.Error(t, err, "Expected an error when getting the task from Redis")
	assert.Equal(t, redis.Nil, err, "Expected redis.Nil error when getting the task from Redis")

	// Verify the task list was updated in Redis
	tasksJSON, err = rdb.Get(context.Background(), tasksKey).Bytes()
	assert.NoError(t, err, "Expected no error when getting the tasks list from Redis")

	var cachedTasks []models.Task
	err = json.Unmarshal([]byte(tasksJSON), &cachedTasks)
	assert.NoError(t, err, "Expected no error when unmarshaling the cached tasks")
	assert.Len(t, cachedTasks, 0, "Expected the tasks list to be empty after deletion")
}


func TestGetUserIdWithUsername(t *testing.T) {
	// Setup the database connection
	dsn := "root:root@tcp(localhost:3306)/to_do"
	db, err := sql.Open("mysql", dsn)
	assert.NoError(t, err, "Failed to connect to the database")
	defer db.Close()

	// Create the repository
	repo := repository.NewRepository(db, nil)

	// Insert a user into the database directly for testing
	username := "testuser"
	_, err = db.Exec("INSERT INTO users (username) VALUES (?)", username)
	assert.NoError(t, err, "Failed to insert user into the database")

	// Call the function to test
	userID, err := repo.GetUserIdWithUsername(username)
	assert.NoError(t, err, "Expected no error from GetUserIdWithUsername")
	assert.NotZero(t, userID, "Expected user ID to be non-zero")

	// Verify that the correct user ID is returned
	var expectedID int64
	err = db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&expectedID)
	assert.NoError(t, err, "Expected no error when querying the user ID from the database")
	assert.Equal(t, expectedID, userID, "Expected user ID to match the inserted user's ID")

	// Test for a non-existing user
	nonExistentUsername := "nonexistentuser"
	userID, err = repo.GetUserIdWithUsername(nonExistentUsername)
	assert.Error(t, err, "Expected an error when querying a non-existing user")
	assert.Zero(t, userID, "Expected user ID to be zero for non-existing user")
}