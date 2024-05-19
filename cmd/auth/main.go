package main

import (
	"database/sql"
	"net"

	"github.com/damirbeybitov/todo_project/internal/auth/repository"
	auth "github.com/damirbeybitov/todo_project/internal/auth/service"
	"github.com/damirbeybitov/todo_project/internal/config"
	"github.com/damirbeybitov/todo_project/internal/log"
	pb "github.com/damirbeybitov/todo_project/proto/auth"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
)

func main() {

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.ErrorLogger.Fatalf("failed to listen: %v", err)
	}

	myConfig, err := config.NewConfig("config.json")
	if err != nil {
		log.ErrorLogger.Fatalf("failed to read config: %v", err)
	}

	db, err := sql.Open("mysql", myConfig.SqlConnection)
	if err != nil {
		log.ErrorLogger.Fatalf("failed to connect to database: %v", err)
	}

	repo := repository.NewRepository(db)

	server := grpc.NewServer()
	authService := auth.NewAuthService(repo) // Создание экземпляра сервиса пользователей
	pb.RegisterAuthServiceServer(server, authService)

	log.InfoLogger.Println("Authentication service is running on port 50051")
	if err := server.Serve(listener); err != nil {
		log.ErrorLogger.Fatalf("failed to serve: %v", err)
	}
}
