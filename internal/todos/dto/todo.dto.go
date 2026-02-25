package dto

import todo_entities "todo-api/internal/todos/entities"

type TodoDto struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

func NewTodoDto(t *todo_entities.Todo) TodoDto {
	return TodoDto{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Done:        t.Done,
	}
}
