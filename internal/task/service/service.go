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
	repo repository.Repository
	taskPB.UnimplementedTaskServiceServer
}

// NewTaskService создает новый экземпляр TaskService.
func NewTaskService(repo repository.Repository) taskPB.TaskServiceServer {
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
		log.ErrorLogger.Printf("Failed to create task: %v", err)
		return nil, err
	}

	return &taskPB.CreateTaskResponse{Id: taskID}, nil
	// return &taskPB.CreateTaskResponse{
	// 	Id: "123",
	// }, nil
}

// GetTask реализует метод получения задачи в рамках интерфейса TaskServiceServer.
func (s *TaskService) GetTask(ctx context.Context, req *taskPB.GetTaskRequest) (*taskPB.GetTaskResponse, error) {
	log.InfoLogger.Printf("Getting task with ID: %s", req.Id)

	// Реализация получения задачи

	// В данном примере просто возвращается фиктивная задача.
	return &taskPB.GetTaskResponse{
		Task: &taskPB.Task{
			Id:          req.Id,
			Title:       "Example Task",
			Description: "This is an example task",
			// Другие поля задачи
		},
	}, nil
}

// UpdateTask реализует метод обновления задачи в рамках интерфейса TaskServiceServer.
func (s *TaskService) UpdateTask(ctx context.Context, req *taskPB.UpdateTaskRequest) (*taskPB.UpdateTaskResponse, error) {
	log.InfoLogger.Printf("Updating task with ID: %s", req.Task.Id)

	// Реализация обновления задачи

	// В данном примере просто возвращается сообщение об успешном обновлении.
	return &taskPB.UpdateTaskResponse{
		Task: req.Task,
	}, nil
}

// DeleteTask реализует метод удаления задачи в рамках интерфейса TaskServiceServer.
func (s *TaskService) DeleteTask(ctx context.Context, req *taskPB.DeleteTaskRequest) (*taskPB.DeleteTaskResponse, error) {
	log.InfoLogger.Printf("Deleting task with ID: %s", req.Id)

	// Реализация удаления задачи

	// В данном примере просто возвращается сообщение об успешном удалении.
	return &taskPB.DeleteTaskResponse{
		Id: req.Id,
	}, nil
}
