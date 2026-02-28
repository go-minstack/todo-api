package todos

import (
	"log/slog"

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
	log   *slog.Logger
}

func NewTodoService(todos *todo_repos.TodoRepository, log *slog.Logger) *TodoService {
	return &TodoService{todos: todos, log: log}
}

func (s *TodoService) List() ([]todo_dto.TodoDto, error) {
	todos, err := s.todos.FindAll()
	if err != nil {
		s.log.Error("failed to list todos", "error", err)
		return nil, err
	}
	s.log.Info("listed todos", "count", len(todos))
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
		s.log.Error("failed to create todo", "error", err)
		return nil, err
	}
	s.log.Info("todo created", "todo_id", todo.ID)
	result := todo_dto.NewTodoDto(todo)
	return &result, nil
}

func (s *TodoService) Get(id uint) (*todo_dto.TodoDto, error) {
	todo, err := s.todos.FindByID(id)
	if err != nil {
		s.log.Error("todo not found", "todo_id", id)
		return nil, err
	}
	result := todo_dto.NewTodoDto(todo)
	return &result, nil
}

func (s *TodoService) Update(id uint, input todo_dto.UpdateTodoDto) (*todo_dto.TodoDto, error) {
	todo, err := s.todos.FindByID(id)
	if err != nil {
		return nil, err
	}

	columns := map[string]interface{}{}
	if input.Title != "" {
		columns["title"] = input.Title
		todo.Title = input.Title
	}
	if input.Description != "" {
		columns["description"] = input.Description
		todo.Description = input.Description
	}
	if input.Done != nil {
		columns["done"] = *input.Done
		todo.Done = *input.Done
	}
	if err := s.todos.UpdatesByID(id, columns); err != nil {
		s.log.Error("failed to update todo", "todo_id", id, "error", err)
		return nil, err
	}

	s.log.Info("todo updated", "todo_id", id)
	result := todo_dto.NewTodoDto(todo)
	return &result, nil
}

func (s *TodoService) Delete(id uint) error {
	if _, err := s.todos.FindByID(id); err != nil {
		s.log.Error("todo not found for deletion", "todo_id", id)
		return err
	}
	if err := s.todos.DeleteByID(id); err != nil {
		s.log.Error("failed to delete todo", "todo_id", id, "error", err)
		return err
	}
	s.log.Info("todo deleted", "todo_id", id)
	return nil
}
