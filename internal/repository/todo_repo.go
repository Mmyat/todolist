package repository

import (
    "todo-api/internal/models"
    "gorm.io/gorm"
)

type TodoRepository struct {
    db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
    return &TodoRepository{db: db}
}

func (r *TodoRepository) Create(todo *models.Todo) error {
    return r.db.Create(todo).Error
}

func (r *TodoRepository) GetAllWithPagination(offset, limit int) ([]models.Todo, error) {
    var todos []models.Todo
    err := r.db.Offset(offset).Limit(limit).Find(&todos).Error
    return todos, err
}