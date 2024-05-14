package main

import (
	"database/sql"
	"net"

	"github.com/damirbeybitov/todo_project/internal/log"
	"github.com/damirbeybitov/todo_project/internal/task"
	pb "github.com/damirbeybitov/todo_project/proto/task"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

func main() {

	listener, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.ErrorLogger.Fatalf("failed to listen: %v", err)
	}

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/to_do")
	if err != nil {
		log.ErrorLogger.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	redis := redis.NewClient(&redis.Options{
        Addr:	  "localhost:6379",
        Password: "", // no password set
        DB:		  0,  // use default DB
    })
	defer redis.Close()

	server := grpc.NewServer()
	taskService := task.NewTaskService(db, redis) 
	pb.RegisterTaskServiceServer(server, taskService)

	log.InfoLogger.Println("Task manager service is running on port 50051")
	if err := server.Serve(listener); err != nil {
		log.ErrorLogger.Fatalf("failed to serve: %v", err)
	}
}
