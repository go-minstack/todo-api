package dto

type UpdateTodoDto struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        *bool  `json:"done"`
}
