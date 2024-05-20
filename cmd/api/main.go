package main

import (
	"log"

	_ "github.com/gin-gonic/gin"
	_ "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"

	"github.com/damirbeybitov/todo_project/internal/handlers"
	"github.com/damirbeybitov/todo_project/internal/models"
	"github.com/damirbeybitov/todo_project/internal/repository"
	"github.com/damirbeybitov/todo_project/internal/service"
	pbAuth "github.com/damirbeybitov/todo_project/proto/auth"
	pbTask "github.com/damirbeybitov/todo_project/proto/task"
	pbUser "github.com/damirbeybitov/todo_project/proto/user"
	"google.golang.org/grpc"
)

// @title Todo Project API
// @version 1.0
// @description API Server for Todo Project

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// Подключение к серверу микросервиса пользователей
	userConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer userConn.Close()

	// Создание клиентского объекта
	userClient := pbUser.NewUserServiceClient(userConn)

	authConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer authConn.Close()

	authClient := pbAuth.NewAuthServiceClient(authConn)

	taskConn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer taskConn.Close()

	taskClient := pbTask.NewTaskServiceClient(taskConn)

	microServiceClients := models.MicroServiceClients{
		UserClient: userClient,
		AuthClient: authClient,
		TaskClient: taskClient,
	}

	repo := repository.NewRepository(microServiceClients)

	handler := handlers.NewHandler(repo)

	service := service.NewService(handler)
	service.LaunchServer()
}
