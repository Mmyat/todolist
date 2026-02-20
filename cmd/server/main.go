package main

import (
    "log"
    "todo-api/configs"
    "todo-api/internal/handler"
    "todo-api/internal/models"
    "todo-api/internal/repository"
    "todo-api/internal/service"
    "github.com/gin-gonic/gin"
)

func main() {
    // Database & Redis Init
    db := configs.GetDB()
    redis := configs.GetRedis()
    // Database Auto-migrate
    db.AutoMigrate(&models.Todo{})
    // Database Connection Pool Setup
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}
	
	// For maximum open connections
	sqlDB.SetMaxOpenConns(100) 
	// For idle connections
	sqlDB.SetMaxIdleConns(50)
    // 2. Dependency Injection
    todoRepo := repository.NewTodoRepository(db)
    todoSvc := service.NewTodoService(todoRepo, redis)
    todoHandler := handler.NewTodoHandler(todoSvc)

    // 3. Routes setup
    r := gin.Default()
    r.GET("/todos", todoHandler.GetTodos)
    r.POST("/todos", todoHandler.CreateTodo)

    r.Run(":8080")
}