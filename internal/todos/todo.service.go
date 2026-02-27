package todos

import (
	"github.com/go-minstack/repository"
	"todo-api/internal/todos/dto"
	todo_entities "todo-api/internal/todos/entities"
	todo_repos "todo-api/internal/todos/repositories"
)

type todoRepository interface {
	FindAll(opts ...repository.QueryOption) ([]todo_entities.Todo, error)
	FindByID(id uint) (*todo_entities.Todo, error)
	Create(entity *todo_entities.Todo) error
	UpdatesByID(id uint, columns map[string]interface{}) error
	DeleteByID(id uint) error
}

type TodoService struct {
	todos todoRepository
}

func NewTodoService(todos *todo_repos.TodoRepository) *TodoService {
	return &TodoService{todos: todos}
}

func (s *TodoService) List() ([]todo_dto.TodoDto, error) {
	todos, err := s.todos.FindAll()
	if err != nil {
		return nil, err
	}
	dtos := make([]todo_dto.TodoDto, len(todos))
	for i, t := range todos {
		dtos[i] = todo_dto.NewTodoDto(&t)
	}
	return dtos, nil
}

func (s *TodoService) Create(input todo_dto.CreateTodoDto) (*todo_dto.TodoDto, error) {
	todo := &todo_entities.Todo{
		Title:       input.Title,
		Description: input.Description,
	}
	if err := s.todos.Create(todo); err != nil {
		return nil, err
	}
	result := todo_dto.NewTodoDto(todo)
	return &result, nil
}

func (s *TodoService) Get(id uint) (*todo_dto.TodoDto, error) {
	todo, err := s.todos.FindByID(id)
	if err != nil {
		return nil, err
	}
	result := todo_dto.NewTodoDto(todo)
	return &result, nil
}

func (s *TodoService) Update(id uint, input todo_dto.UpdateTodoDto) (*todo_dto.TodoDto, error) {
	if _, err := s.todos.FindByID(id); err != nil {
		return nil, err
	}

	columns := map[string]interface{}{}
	if input.Title != "" {
		columns["title"] = input.Title
	}
	if input.Description != "" {
		columns["description"] = input.Description
	}
	if input.Done != nil {
		columns["done"] = *input.Done
	}
	if err := s.todos.UpdatesByID(id, columns); err != nil {
		return nil, err
	}

	return s.Get(id)
}

func (s *TodoService) Delete(id uint) error {
	if _, err := s.todos.FindByID(id); err != nil {
		return err
	}
	return s.todos.DeleteByID(id)
}
