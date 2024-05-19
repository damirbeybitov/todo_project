package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/damirbeybitov/todo_project/internal/log"
	"github.com/damirbeybitov/todo_project/internal/models"
	"github.com/damirbeybitov/todo_project/internal/repository"

	pbAuth "github.com/damirbeybitov/todo_project/proto/auth"
	// pbTask "github.com/damirbeybitov/todo_project/proto/task"
	pbUser "github.com/damirbeybitov/todo_project/proto/user"
)

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

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