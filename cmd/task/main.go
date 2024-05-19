package main

import (
	"database/sql"
	"net"

	"github.com/damirbeybitov/todo_project/internal/config"
	"github.com/damirbeybitov/todo_project/internal/log"
	"github.com/damirbeybitov/todo_project/internal/task/repository"
	task "github.com/damirbeybitov/todo_project/internal/task/service"
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

	myConfig, err := config.NewConfig("config.json")
	if err != nil {
		log.ErrorLogger.Fatalf("failed to read config: %v", err)
	}

	db, err := sql.Open("mysql", myConfig.SqlConnection)
	if err != nil {
		log.ErrorLogger.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer redisClient.Close()

	repo := repository.NewRepository(db)

	server := grpc.NewServer()
	taskService := task.NewTaskService(repo, redisClient)
	pb.RegisterTaskServiceServer(server, taskService)

	log.InfoLogger.Println("Task manager service is running on port 50053")
	if err := server.Serve(listener); err != nil {
		log.ErrorLogger.Fatalf("failed to serve: %v", err)
	}
}
