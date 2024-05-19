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

	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", s.handler.RegisterHandler).Methods("POST")
	authRouter.HandleFunc("/login", s.handler.LoginHandler).Methods("GET")
	authRouter.HandleFunc("/refresh_token", s.handler.RefreshTokenHandler).Methods("GET")

	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.Use(s.handler.UserIdentity)
	userRouter.HandleFunc("/get_user_profile", s.handler.GetUserProfileHandler).Methods("GET")
	userRouter.HandleFunc("/delete_user", s.handler.DeleteUserHandler).Methods("DELETE")

	taskRouter := router.PathPrefix("/task").Subrouter()
	taskRouter.Use(s.handler.UserIdentity)
	taskRouter.HandleFunc("/create_task", s.handler.CreateTaskHandler).Methods("POST")
	taskRouter.HandleFunc("/get_tasks", s.handler.GetTasksHandler).Methods("GET")
	taskRouter.HandleFunc("/get_task/{id}", s.handler.GetTaskHandler).Methods("GET")
	taskRouter.HandleFunc("/update_task", s.handler.UpdateTaskHandler).Methods("PUT")
	taskRouter.HandleFunc("/delete_task/{id}", s.handler.DeleteTaskHandler).Methods("DELETE")

	log.InfoLogger.Print("Main service is running on port 8080")
	http.ListenAndServe(":8080", router)
}
