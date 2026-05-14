package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"todo-api/internal/models"
	"todo-api/internal/repository"

	"github.com/redis/go-redis/v9"
)

type TodoService struct {
	repo  *repository.TodoRepository
	redis *redis.Client
}

func NewTodoService(repo *repository.TodoRepository, redis *redis.Client) *TodoService {
	return &TodoService{repo: repo, redis: redis}
}

func (s *TodoService) CreateTodo(ctx context.Context, todo *models.Todo) error {
	if todo.Title == "" {
		return errors.New("title cannot be empty")
	}
	err := s.repo.Create(todo)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	
	// Cache invalidation using Redis Sets
	go func() {
		bgCtx := context.Background()
		keys, err := s.redis.SMembers(bgCtx, "todos:pagination_keys").Result()
		if err == nil && len(keys) > 0 {
			// Delete all cached pages
			s.redis.Del(bgCtx, keys...)
			// Clear the tracking set
			s.redis.Del(bgCtx, "todos:pagination_keys")
		} else if err != nil && err != redis.Nil {
			log.Printf("Failed to invalidate cache: %v\n", err)
		}
	}()
	return nil
}

func (s *TodoService) GetAllTodosWithPagination(ctx context.Context, page, limit int) ([]models.Todo, error) {
	offset := (page - 1) * limit
	cacheKey := fmt.Sprintf("todos:p:%d:l:%d", page, limit)
	
	val, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var todos []models.Todo
		if err := json.Unmarshal([]byte(val), &todos); err == nil {
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
		// Track this cache key in our set
		s.redis.SAdd(ctx, "todos:pagination_keys", cacheKey)
	}
	return todos, nil
}