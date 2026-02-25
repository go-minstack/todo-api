package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-minstack/core"
	mgin "github.com/go-minstack/gin"
	"github.com/go-minstack/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"todo-api/internal/todos"
	todo_entities "todo-api/internal/todos/entities"
	"gorm.io/gorm"
)

func setupApp(t *testing.T) *gin.Engine {
	t.Helper()

	t.Setenv("MINSTACK_DB_URL", ":memory:")
	t.Setenv("MINSTACK_PORT", "0")
	gin.SetMode(gin.TestMode)

	app := core.New(mgin.Module(), sqlite.Module())
	todos.Register(app)

	var engine *gin.Engine
	app.Invoke(func(r *gin.Engine) { engine = r })
	app.Invoke(func(db *gorm.DB) error {
		return db.AutoMigrate(&todo_entities.Todo{})
	})

	ctx := context.Background()
	require.NoError(t, app.Start(ctx))
	t.Cleanup(func() { app.Stop(ctx) })

	return engine
}

func jsonBody(data any) *bytes.Buffer {
	b, _ := json.Marshal(data)
	return bytes.NewBuffer(b)
}

func parseJSON(t *testing.T, w *httptest.ResponseRecorder, v any) {
	t.Helper()
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), v))
}

func TestTodoAPI(t *testing.T) {
	r := setupApp(t)

	var todoID uint

	t.Run("List empty", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/todos", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var todos []map[string]any
		parseJSON(t, w, &todos)
		assert.Empty(t, todos)
	})

	t.Run("Create", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := jsonBody(map[string]string{
			"title":       "Buy milk",
			"description": "From the store",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/todos", body)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var todo map[string]any
		parseJSON(t, w, &todo)
		assert.Equal(t, "Buy milk", todo["title"])
		assert.Equal(t, "From the store", todo["description"])
		assert.False(t, todo["done"].(bool))
		todoID = uint(todo["id"].(float64))
		assert.NotZero(t, todoID)
	})

	t.Run("List with todo", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/todos", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var todos []map[string]any
		parseJSON(t, w, &todos)
		assert.Len(t, todos, 1)
		assert.Equal(t, "Buy milk", todos[0]["title"])
	})

	t.Run("Get", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/todos/%d", todoID), nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var todo map[string]any
		parseJSON(t, w, &todo)
		assert.Equal(t, "Buy milk", todo["title"])
	})

	t.Run("Update", func(t *testing.T) {
		done := true
		w := httptest.NewRecorder()
		body := jsonBody(map[string]any{
			"title": "Buy bread",
			"done":  done,
		})
		req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/todos/%d", todoID), body)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var todo map[string]any
		parseJSON(t, w, &todo)
		assert.Equal(t, "Buy bread", todo["title"])
		assert.True(t, todo["done"].(bool))
	})

	t.Run("Delete", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/todos/%d", todoID), nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("Get deleted", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/todos/%d", todoID), nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Create without title", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := jsonBody(map[string]string{"description": "No title"})
		req := httptest.NewRequest(http.MethodPost, "/api/todos", body)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
