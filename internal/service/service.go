package service

import (
	"net/http"

	"github.com/damirbeybitov/todo_project/internal/handlers"
	"github.com/damirbeybitov/todo_project/internal/log"
	"github.com/gorilla/mux"
)

type Service struct {
	handler *handlers.Handler
}

func NewService(handler *handlers.Handler) *Service {
	return &Service{handler: handler}
}

func (s *Service) LaunchServer() {
	router := mux.NewRouter()
	router.HandleFunc("/register", s.handler.RegisterHandler).Methods("POST")
	router.HandleFunc("/login", s.handler.LoginHandler).Methods("GET")
	// router.HandleFunc("/get_user_profile", s.getUserProfile).Methods("GET")
	// router.HandleFunc("/create_task", s.createTask).Methods("POST")
	// router.HandleFunc("/get_tasks", s.getTasks).Methods("GET")
	// router.HandleFunc("/get_task", s.getTask).Methods("GET")
	// router.HandleFunc("/update_task", s.updateTask).Methods("PUT")
	// router.HandleFunc("/delete_task", s.deleteTask).Methods("DELETE")

	http.ListenAndServe(":8080", router)
	log.InfoLogger.Print("Authentication service is running on port 8080")
}
