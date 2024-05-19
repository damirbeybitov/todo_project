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
	// router.HandleFunc("/create_task", s.createTask).Methods("POST")
	router.HandleFunc("/get_tasks/{id}", s.handler.GetTasksHandler).Methods("GET")
	// router.HandleFunc("/get_task", s.getTask).Methods("GET")
	// router.HandleFunc("/update_task", s.updateTask).Methods("PUT")
	// router.HandleFunc("/delete_task", s.deleteTask).Methods("DELETE")

	log.InfoLogger.Print("Main service is running on port 8080")
	http.ListenAndServe(":8080", router)
}
