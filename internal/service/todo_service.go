package service

import (
    "context"
	 "strconv"
    "encoding/json"
    "time"
	 "errors"
	 "fmt"
    "todo-api/internal/models"
	 "todo-api/internal/repository"
    "github.com/redis/go-redis/v9"
)
type TodoService struct {
    repo *repository.TodoRepository
    redis *redis.Client
}

func NewTodoService(repo *repository.TodoRepository, redis *redis.Client) *TodoService {
    return &TodoService{repo: repo, redis: redis}
}

func (s *TodoService) CreateTodo(ctx context.Context,todo *models.Todo) error {
    if todo.Title == "" {
        return errors.New("title cannot be empty")
    }
	err := s.repo.Create(todo)
	if err != nil {
        return fmt.Errorf("database error: %w", err)
    }
	go func() {
        cacheKey := "all_todos"
        s.redis.Del(context.Background(), cacheKey)
        fmt.Println("Redis cache cleared for data consistency")
   }()
	return nil
}

func (s *TodoService) GetAllTodosWithPagination(ctx context.Context, page, limit int) ([]models.Todo, error) {
	offset := (page - 1) * limit
	cacheKey := fmt.Sprintf("todos:p:%d:l:%d", page, limit)
	val, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var todos []models.Todo
		start := time.Now()
		if err := json.Unmarshal([]byte(val), &todos); err == nil {
			println("redis is working for getting")
			fmt.Printf("Redis Get took: %v\n", time.Since(start))
			return todos, nil
		}
	}

	todos, err := s.repo.GetAllWithPagination(offset, limit)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(todos)
	if err == nil {
		s.redis.Set(ctx, cacheKey, data, 10*time.Minute)
	}
	go func() {
        data, _ := json.Marshal(todos)
        s.redis.Set(context.Background(), cacheKey, data, 10*time.Minute)
   }()
	return todos, nil
}