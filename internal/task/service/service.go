package task

import (
	"context"

	"github.com/damirbeybitov/todo_project/internal/log"
	"github.com/damirbeybitov/todo_project/internal/models"
	"github.com/damirbeybitov/todo_project/internal/task/repository"
	taskPB "github.com/damirbeybitov/todo_project/proto/task"
)

// TaskService представляет сервис управления задачами.
type TaskService struct {
	repo *repository.Repository
	taskPB.UnimplementedTaskServiceServer
}

// NewTaskService создает новый экземпляр TaskService.
func NewTaskService(repo *repository.Repository) taskPB.TaskServiceServer {
	return &TaskService{repo: repo}
}

// CreateTask реализует метод создания задачи в рамках интерфейса TaskServiceServer.
func (s *TaskService) CreateTask(ctx context.Context, req *taskPB.CreateTaskRequest) (*taskPB.CreateTaskResponse, error) {
	log.InfoLogger.Printf("Creating task with title: %s", req.Task.Title)

	// Реализация создания задачи

	// В данном примере просто возвращается фиктивный идентификатор задачи.
	task := models.Task{
		Title:       req.Task.Title,
		Description: req.Task.Description,
		Status:      req.Task.Status,
		UserID:      req.Task.UserId,
	}

	taskID, err := s.repo.CreateTask(task)
	if err != nil {
		return nil, err
	}

	return &taskPB.CreateTaskResponse{Id: taskID}, nil
}

// GetTask реализует метод получения задачи в рамках интерфейса TaskServiceServer.
func (s *TaskService) GetTask(ctx context.Context, req *taskPB.GetTaskRequest) (*taskPB.GetTaskResponse, error) {
	log.InfoLogger.Printf("Getting task with ID: %d", req.Id)

	// Реализация получения задачи
	task, err := s.repo.GetTaskByID(req.Id)
	if err != nil {
		return nil, err
	}

	log.InfoLogger.Printf("Task found: %v", task)
	// В данном примере просто возвращается фиктивная задача.
	return &taskPB.GetTaskResponse{
		Task: &taskPB.Task{
			Id:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			UserId:      task.UserID,
		},
	}, nil
}

// UpdateTask реализует метод обновления задачи в рамках интерфейса TaskServiceServer.
func (s *TaskService) UpdateTask(ctx context.Context, req *taskPB.UpdateTaskRequest) (*taskPB.UpdateTaskResponse, error) {
	log.InfoLogger.Printf("Updating task with ID: %s", req.Task.Id)

	// Реализация обновления задачи
	task := models.Task{
		ID:          req.Task.Id,
		Title:       req.Task.Title,
		Description: req.Task.Description,
		Status:      req.Task.Status,
		UserID:      req.Task.UserId,
	}

	err := s.repo.UpdateTask(task)
	if err != nil {
		return nil, err
	}

	log.InfoLogger.Printf("Task updated: %v", task)
	// В данном примере просто возвращается сообщение об успешном обновлении.
	return &taskPB.UpdateTaskResponse{
		Task: &taskPB.Task{
			Id:          req.Task.Id,
			Title:       req.Task.Title,
			Description: req.Task.Description,
			Status:      req.Task.Status,
			UserId:      req.Task.UserId,
		},
	}, nil
}

// DeleteTask реализует метод удаления задачи в рамках интерфейса TaskServiceServer.
func (s *TaskService) DeleteTask(ctx context.Context, req *taskPB.DeleteTaskRequest) (*taskPB.DeleteTaskResponse, error) {
	log.InfoLogger.Printf("Deleting task with ID: %s", req.Id)

	// Реализация удаления задачи
	err := s.repo.DeleteTask(req.Id)
	if err != nil {
		return nil, err
	}

	log.InfoLogger.Printf("Task deleted with ID: %s", req.Id)
	// В данном примере просто возвращается сообщение об успешном удалении.
	return &taskPB.DeleteTaskResponse{Id: req.Id}, nil
}
