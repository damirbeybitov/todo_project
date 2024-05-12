package main

import (
	"log"
	"net"

	"github.com/damirbeybitov/todo_project/internal/user"
	pb "github.com/damirbeybitov/todo_project/proto/user"
	"google.golang.org/grpc"
)

func main() {

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	userService := user.NewUserService() // Создание экземпляра сервиса пользователей
	pb.RegisterUserServiceServer(server, userService)

	log.Println("User service is running on port 50051")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
