package main

import (
	"context"
	"dq_scheduler_v2/config"
	"dq_scheduler_v2/executor"
	"dq_scheduler_v2/handler"
	"dq_scheduler_v2/service"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// MongoDB connection URI
	mongoURI := "mongodb://localhost:27017"

	// Create MongoDB client
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Create config
	config, err := config.NewConfig(mongoURI)
	if err != nil {
		log.Fatalf("Failed to create config: %v", err)
	}

	// Create job executor
	jobExecutor := executor.NewJobExecutor()

	// Create scheduler service
	schedulerService := service.NewSchedulerService(config, jobExecutor.ExecuteJob)

	// Create scheduler handler
	schedulerHandler := handler.NewSchedulerHandler(schedulerService, config)

	// Start scheduled jobs
	schedulerService.StartScheduler()

	// Set up Gin router
	router := gin.Default()

	// Register API routes
	router.POST("/schedulers", schedulerHandler.CreateScheduler)
	router.GET("/schedulers", schedulerHandler.ListSchedulers)
	router.GET("/schedulers/:id", schedulerHandler.GetScheduler)
	router.PUT("/schedulers/:id", schedulerHandler.UpdateScheduler)
	router.DELETE("/schedulers/:id", schedulerHandler.DeleteScheduler)
	router.POST("/schedulers/:id/start", schedulerHandler.StartScheduler)
	router.POST("/schedulers/:id/stop", schedulerHandler.StopScheduler)
	router.POST("/schedulers/:id/resume", schedulerHandler.ResumeScheduler)

	// Start HTTP server
	log.Fatal(router.Run(":8080"))
}
