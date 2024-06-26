package main

import (
	"database/sql"
	"net"

	"github.com/damirbeybitov/todo_project/internal/config"
	"github.com/damirbeybitov/todo_project/internal/log"
	"github.com/damirbeybitov/todo_project/internal/user/repository"
	user "github.com/damirbeybitov/todo_project/internal/user/serivice"
	pb "github.com/damirbeybitov/todo_project/proto/user"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
)

func main() {

	listener, err := net.Listen("tcp", ":50051")
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
	userService := user.NewUserService(repo) // Создание экземпляра сервиса пользователей
	pb.RegisterUserServiceServer(server, userService)

	log.InfoLogger.Println("User service is running on port 50051")
	if err := server.Serve(listener); err != nil {
		log.ErrorLogger.Fatalf("failed to serve: %v", err)
	}
}
