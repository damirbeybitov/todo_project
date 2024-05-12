package main

import (
	"database/sql"
	"net"

	"github.com/damirbeybitov/todo_project/internal/log"
	"github.com/damirbeybitov/todo_project/internal/user"
	pb "github.com/damirbeybitov/todo_project/proto/user"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
)

func main() {

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.ErrorLogger.Fatalf("failed to listen: %v", err)
	}

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/to_do")
	if err != nil {
		log.ErrorLogger.Fatalf("failed to connect to database: %v", err)
	}

	server := grpc.NewServer()
	userService := user.NewUserService(db) // Создание экземпляра сервиса пользователей
	pb.RegisterUserServiceServer(server, userService)

	log.InfoLogger.Println("User service is running on port 50051")
	if err := server.Serve(listener); err != nil {
		log.ErrorLogger.Fatalf("failed to serve: %v", err)
	}
}
