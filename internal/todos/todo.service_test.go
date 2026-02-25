package todos

import (
	"errors"
	"testing"

	"github.com/go-minstack/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"todo-api/internal/todos/dto"
	todo_entities "todo-api/internal/todos/entities"
)

// mockTodoRepo is an in-memory implementation of todoRepository for unit tests.
type mockTodoRepo struct {
	todos  []todo_entities.Todo
	nextID uint
}

func newMockTodoRepo() *mockTodoRepo {
	return &mockTodoRepo{nextID: 1}
}

func (m *mockTodoRepo) FindAll(opts ...repository.QueryOption) ([]todo_entities.Todo, error) {
	result := make([]todo_entities.Todo, len(m.todos))
	copy(result, m.todos)
	return result, nil
}

func (m *mockTodoRepo) FindByID(id uint) (*todo_entities.Todo, error) {
	for _, t := range m.todos {
		if t.ID == id {
			return &t, nil
		}
	}
	return nil, errors.New("record not found")
}

func (m *mockTodoRepo) Create(entity *todo_entities.Todo) error {
	entity.ID = m.nextID
	m.nextID++
	m.todos = append(m.todos, *entity)
	return nil
}

func (m *mockTodoRepo) UpdatesByID(id uint, columns map[string]interface{}) error {
	for i, t := range m.todos {
		if t.ID == id {
			if v, ok := columns["title"]; ok {
				m.todos[i].Title = v.(string)
			}
			if v, ok := columns["description"]; ok {
				m.todos[i].Description = v.(string)
			}
			if v, ok := columns["done"]; ok {
				m.todos[i].Done = v.(bool)
			}
			return nil
		}
	}
	return errors.New("record not found")
}

func (m *mockTodoRepo) DeleteByID(id uint) error {
	for i, t := range m.todos {
		if t.ID == id {
			m.todos = append(m.todos[:i], m.todos[i+1:]...)
			return nil
		}
	}
	return errors.New("record not found")
}

func TestTodoService_List(t *testing.T) {
	repo := newMockTodoRepo()
	svc := &TodoService{todos: repo}

	// empty list
	todos, err := svc.List()
	require.NoError(t, err)
	assert.Empty(t, todos)

	// add some todos
	repo.Create(&todo_entities.Todo{Title: "Buy milk"})
	repo.Create(&todo_entities.Todo{Title: "Walk dog"})

	todos, err = svc.List()
	require.NoError(t, err)
	assert.Len(t, todos, 2)
	assert.Equal(t, "Buy milk", todos[0].Title)
	assert.Equal(t, "Walk dog", todos[1].Title)
}

func TestTodoService_Create(t *testing.T) {
	repo := newMockTodoRepo()
	svc := &TodoService{todos: repo}

	todo, err := svc.Create(dto.CreateTodoDto{
		Title:       "Buy milk",
		Description: "From the store",
	})
	require.NoError(t, err)
	assert.Equal(t, uint(1), todo.ID)
	assert.Equal(t, "Buy milk", todo.Title)
	assert.Equal(t, "From the store", todo.Description)
	assert.False(t, todo.Done)

	// verify it was persisted
	assert.Len(t, repo.todos, 1)
}

func TestTodoService_Get(t *testing.T) {
	repo := newMockTodoRepo()
	svc := &TodoService{todos: repo}

	repo.Create(&todo_entities.Todo{Title: "Buy milk"})

	t.Run("found", func(t *testing.T) {
		todo, err := svc.Get(1)
		require.NoError(t, err)
		assert.Equal(t, "Buy milk", todo.Title)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := svc.Get(999)
		assert.Error(t, err)
	})
}

func TestTodoService_Update(t *testing.T) {
	repo := newMockTodoRepo()
	svc := &TodoService{todos: repo}

	repo.Create(&todo_entities.Todo{Title: "Buy milk", Description: "From the store"})

	t.Run("update title", func(t *testing.T) {
		todo, err := svc.Update(1, dto.UpdateTodoDto{Title: "Buy bread"})
		require.NoError(t, err)
		assert.Equal(t, "Buy bread", todo.Title)
	})

	t.Run("mark done", func(t *testing.T) {
		done := true
		todo, err := svc.Update(1, dto.UpdateTodoDto{Done: &done})
		require.NoError(t, err)
		assert.True(t, todo.Done)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := svc.Update(999, dto.UpdateTodoDto{Title: "Nope"})
		assert.Error(t, err)
	})
}

func TestTodoService_Delete(t *testing.T) {
	repo := newMockTodoRepo()
	svc := &TodoService{todos: repo}

	repo.Create(&todo_entities.Todo{Title: "Buy milk"})

	t.Run("delete existing", func(t *testing.T) {
		err := svc.Delete(1)
		require.NoError(t, err)
		assert.Empty(t, repo.todos)
	})

	t.Run("delete non-existent", func(t *testing.T) {
		err := svc.Delete(999)
		assert.Error(t, err)
	})
}
