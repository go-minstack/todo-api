package todo_repositories

import (
	"github.com/go-minstack/repository"
	todo_entities "todo-api/internal/todos/entities"
	"gorm.io/gorm"
)

type TodoRepository struct {
	*repository.Repository[todo_entities.Todo]
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{repository.NewRepository[todo_entities.Todo](db)}
}
