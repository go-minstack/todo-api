package main

import (
	"github.com/go-minstack/go-minstack/core"
	mgin "github.com/go-minstack/go-minstack/gin"
	"github.com/go-minstack/go-minstack/sqlite"
	"gorm.io/gorm"
	"todo-api/internal/todos"
	todo_entities "todo-api/internal/todos/entities"
)

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(&todo_entities.Todo{})
}

func main() {
	app := core.New(mgin.Module(), sqlite.Module())

	todos.Register(app)

	app.Invoke(migrate)
	app.Run()
}
