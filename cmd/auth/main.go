package main

import (
	"database/sql"
	"net"

	"github.com/damirbeybitov/todo_project/internal/auth"
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

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/to_do")
	if err != nil {
		log.ErrorLogger.Fatalf("failed to connect to database: %v", err)
	}

	server := grpc.NewServer()
	authService := auth.NewAuthService(db) // Создание экземпляра сервиса пользователей
	pb.RegisterAuthServiceServer(server, authService)

	log.InfoLogger.Println("Authentication service is running on port 50051")
	if err := server.Serve(listener); err != nil {
		log.ErrorLogger.Fatalf("failed to serve: %v", err)
	}
}
