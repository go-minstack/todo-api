package todo_dto

type CreateTodoDto struct {
	Title       string `json:"title"       binding:"required"`
	Description string `json:"description"`
}
