package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/damirbeybitov/todo_project/internal/log"
	"github.com/damirbeybitov/todo_project/internal/models"
	"github.com/damirbeybitov/todo_project/internal/repository"
	"github.com/gorilla/mux"

	pbAuth "github.com/damirbeybitov/todo_project/proto/auth"
	pbTask "github.com/damirbeybitov/todo_project/proto/task"
	pbUser "github.com/damirbeybitov/todo_project/proto/user"
)

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

// @Summary Register user
// @Tags auth
// @Description register a new user
// @ID register-user
// @Accept json
// @Produce json
// @Param body body models.RegisterRequest true "User registration data"
// @Success 200 {object} models.RegisterResponse
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /auth/register [post]
func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.ErrorLogger.Printf("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Email == "" || user.Password == "" {
		log.ErrorLogger.Print("Missing required fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	pbResponse, err := h.repo.MicroServiceClients.UserClient.RegisterUser(r.Context(), &pbUser.RegisterUserRequest{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	response := models.RegisterResponse{
		Id: pbResponse.Id,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.ErrorLogger.Printf("Failed to marshal response: %v", err)
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
	log.InfoLogger.Print("Register endpoint done successfully")
}

// @Summary User login
// @Tags auth
// @Description login for existing user
// @ID login-user
// @Accept json
// @Produce json
// @Param body body models.LoginRequest true "User login data"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /auth/login [post]
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.ErrorLogger.Printf("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if user.Username == "" || user.Password == "" {
		log.ErrorLogger.Print("Missing required fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	pbResponse, err := h.repo.MicroServiceClients.AuthClient.Authenticate(r.Context(), &pbAuth.AuthenticateRequest{
		Username: user.Username,
		Password: user.Password,
	})
	if err != nil {
		http.Error(w, "Failed to login user", http.StatusInternalServerError)
		return
	}

	response := models.LoginResponse{
		AccessToken:  pbResponse.AccessToken,
		RefreshToken: pbResponse.RefreshToken,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.ErrorLogger.Printf("Failed to marshal response: %v", err)
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
	log.InfoLogger.Print("Login endpoint done successfully")
}

// @Summary Refresh user token
// @Tags auth
// @Description refresh user access token
// @ID refresh-token
// @Accept json
// @Produce json
// @Param body body models.RefreshTokenRequest true "User refresh token data"
// @Success 200 {object} models.RefreshTokenResponse
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /auth/refresh [post]
func (h *Handler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var refreshToken models.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&refreshToken); err != nil {
		log.ErrorLogger.Printf("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if refreshToken.RefreshToken == "" {
		log.ErrorLogger.Print("Missing required fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	accessToken, err := h.repo.MicroServiceClients.AuthClient.RefreshToken(context.Background(), &pbAuth.RefreshTokenRequest{
		RefreshToken: refreshToken.RefreshToken,
	})
	if err != nil {
		http.Error(w, "Failed to refresh token", http.StatusInternalServerError)
		return
	}

	response := models.RefreshTokenResponse{
		AccessToken: accessToken.AccessToken,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.ErrorLogger.Printf("Failed to marshal response: %v", err)
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
	log.InfoLogger.Print("Refresh token endpoint done successfully")
}

// @Summary Get user profile
// @Tags user
// @Description get user profile by token
// @ID get-user-profile
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.GetUserProfileResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /user/profile [get]
func (h *Handler) GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	userToken := r.Header.Get("Authorization")
	if userToken == "" {
		log.ErrorLogger.Print("Token is missing in Header")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := h.repo.GetUserIdFromRequest(userToken)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	pbResponse, err := h.repo.MicroServiceClients.UserClient.GetUserProfile(r.Context(), &pbUser.GetUserProfileRequest{
		Id: id,
	})
	if err != nil {
		http.Error(w, "Failed to get user profile", http.StatusInternalServerError)
		return
	}

	response := models.GetUserProfileResponse{
		Id:       pbResponse.User.Id,
		Username: pbResponse.User.Username,
		Email:    pbResponse.User.Email,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.ErrorLogger.Printf("Failed to marshal response: %v", err)
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
	log.InfoLogger.Print("Get user profile endpoint done successfully")

}

// @Summary Delete user
// @Tags user
// @Description delete user by username and password
// @ID delete-user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body models.DeleteUserRequest true "User deletion data"
// @Success 200 {object} models.DeleteUserResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /user/delete [post]
func (h *Handler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.DeleteUserRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.ErrorLogger.Printf("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userToken := r.Header.Get("Authorization")
	if userToken == "" {
		log.ErrorLogger.Print("Token is missing in Header")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	username, err := h.repo.GetUsernameFromRequest(userToken)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if user.Password == "" {
		log.ErrorLogger.Print("Missing required fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	log.InfoLogger.Printf("Username to delete user: %s", username)
	pbResponse, err := h.repo.MicroServiceClients.UserClient.DeleteUser(r.Context(), &pbUser.DeleteUserRequest{
		Username: username,
		Password: user.Password,
	})
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	response := models.DeleteUserResponse{
		Message: pbResponse.Message,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.ErrorLogger.Printf("Failed to marshal response: %v", err)
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
	log.InfoLogger.Print("Delete user endpoint done successfully")
}

// @Summary Create task
// @Tags task
// @Description create a new task
// @ID create-task
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body models.CreateTaskRequest true "Task creation data"
// @Success 200 {object} models.CreateTaskResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /task/create [post]
func (h *Handler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.ErrorLogger.Printf("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.Description == "" || req.UserId == 0 {
		log.ErrorLogger.Print("Missing required fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	task := models.Task{
		Title:       req.Title,
		Description: req.Description,
		UserId:      req.UserId,
	}

	taskID, err := h.repo.MicroServiceClients.TaskClient.CreateTask(r.Context(), &pbTask.CreateTaskRequest{
		Task: &pbTask.Task{
			Title:       task.Title,
			Description: task.Description,
			UserId:      task.UserId,
		},
	})

	if err != nil {
		log.ErrorLogger.Printf("Failed to create task: %v", err)
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	response := models.CreateTaskResponse{
		Id: taskID.Id,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.ErrorLogger.Printf("Failed to marshal response: %v", err)
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
	log.InfoLogger.Print("Create task endpoint done successfully")
}

// @Summary Get tasks
// @Tags task
// @Description retrieve all tasks
// @ID get-tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.GetTasksResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /task/all [get]
func (h *Handler) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	// Implement logic to retrieve all tasks
	// Placeholder implementation

	userToken := r.Header.Get("Authorization")
	if userToken == "" {
		log.ErrorLogger.Print("Token is missing in Header")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	username, err := h.repo.GetUsernameFromRequest(userToken)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	pbTasks, err := h.repo.MicroServiceClients.TaskClient.GetTasks(r.Context(), &pbTask.GetTasksRequest{
		Username: username,
	})
	if err != nil {
		http.Error(w, "Failed to get tasks", http.StatusInternalServerError)
		return
	}

	var tasks []models.Task
	pbTask := pbTasks.GetTasks()
	for _, task := range pbTask {
		tasks = append(tasks, models.Task{
			Id:          task.Id,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			UserId:      task.UserId,
		})
	}

	response := models.GetTasksResponse{
		Tasks: tasks,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.ErrorLogger.Printf("Failed to marshal response: %v", err)
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
	log.InfoLogger.Print("Get tasks endpoint done successfully")
}

// @Summary Get task by ID
// @Tags task
// @Description retrieve a task by its ID
// @ID get-task-by-id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Task ID"
// @Success 200 {object} models.Task
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /task/{id} [get]
func (h *Handler) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		log.ErrorLogger.Print("Missing task ID")
		http.Error(w, "Missing task ID", http.StatusBadRequest)
		return
	}

	taskID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.ErrorLogger.Printf("Invalid task ID: %v", err)
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := h.repo.MicroServiceClients.TaskClient.GetTask(r.Context(), &pbTask.GetTaskRequest{
		Id: taskID,
	})
	if err != nil {
		log.ErrorLogger.Printf("Failed to get task: %v", err)
		http.Error(w, "Failed to get task", http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(task)
	if err != nil {
		log.ErrorLogger.Printf("Failed to marshal response: %v", err)
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
	log.InfoLogger.Print("Get task endpoint done successfully")
}

// @Summary Update task
// @Tags task
// @Description update a task
// @ID update-task
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body models.UpdateTaskRequest true "Task update data"
// @Success 200 {object} models.UpdateTaskResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /task/update [post]
func (h *Handler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.ErrorLogger.Printf("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userToken := r.Header.Get("Authorization")
	if userToken == "" {
		log.ErrorLogger.Print("Token is missing in Header")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userId, err := h.repo.GetUserIdFromRequest(userToken)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if req.Id == 0 || req.Title == "" || req.Description == "" || userId == 0 {
		log.ErrorLogger.Print("Missing required fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	task := models.Task{
		Id:          req.Id,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		UserId:      userId,
	}

	UpdateTaskResponse, err := h.repo.MicroServiceClients.TaskClient.UpdateTask(r.Context(), &pbTask.UpdateTaskRequest{
		Task: &pbTask.Task{
			Id:          task.Id,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			UserId:      userId,
		},
	})

	if err != nil {
		log.ErrorLogger.Printf("Failed to update task: %v", err)
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	response := models.UpdateTaskResponse{
		Id:          UpdateTaskResponse.Task.Id,
		Title:       UpdateTaskResponse.Task.Title,
		Description: UpdateTaskResponse.Task.Description,
		Status:      UpdateTaskResponse.Task.Status,
		UserId:      UpdateTaskResponse.Task.UserId,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.ErrorLogger.Printf("Failed to marshal response: %v", err)
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
	log.InfoLogger.Print("Update task endpoint done successfully")
}

// @Summary Delete task
// @Tags task
// @Description delete a task by ID
// @ID delete-task
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Task ID"
// @Success 200 {object} models.DeleteTaskResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /task/{id} [delete]
func (h *Handler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		log.ErrorLogger.Print("Missing task ID")
		http.Error(w, "Missing task ID", http.StatusBadRequest)
		return
	}

	taskID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.ErrorLogger.Printf("Invalid task ID: %v", err)
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	DeleteTaskResponse, err := h.repo.MicroServiceClients.TaskClient.DeleteTask(r.Context(), &pbTask.DeleteTaskRequest{
		Id: taskID,
	})
	if err != nil {
		log.ErrorLogger.Printf("Failed to delete task: %v", err)
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	response := models.DeleteTaskResponse{
		Message: DeleteTaskResponse.Message,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.ErrorLogger.Printf("Failed to marshal response: %v", err)
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
	log.InfoLogger.Print("Delete task endpoint done successfully")
}
