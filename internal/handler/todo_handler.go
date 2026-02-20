package handler

import (
    "net/http"
    "strconv"
    "todo-api/internal/models"
    "todo-api/internal/service"
    "github.com/gin-gonic/gin"
)

type TodoHandler struct {
    service *service.TodoService
}

func NewTodoHandler(service *service.TodoService) *TodoHandler {
    return &TodoHandler{service: service}
}

func (h *TodoHandler) GetTodos(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))    
    todos, err := h.service.GetAllTodosWithPagination(c.Request.Context(), page, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todos"})
        return
    }
    c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {
    var todo models.Todo
    if err := c.ShouldBindJSON(&todo); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Create အောင်မြင်၊ မအောင်မြင်ကို စစ်ဆေးပြီးမှ Status ပြန်ပေးပါ
    if err := h.service.CreateTodo(c.Request.Context(), &todo); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, todo)
}