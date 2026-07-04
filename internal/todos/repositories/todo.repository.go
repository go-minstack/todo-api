package todo_repositories

import (
	"github.com/go-minstack/go-minstack/repository"
	"gorm.io/gorm"
	todo_entities "todo-api/internal/todos/entities"
)

type TodoRepository struct {
	*repository.Repository[todo_entities.Todo]
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{repository.NewRepository[todo_entities.Todo](db)}
}
