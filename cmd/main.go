package main

import (
	"github.com/go-minstack/core"
	mgin "github.com/go-minstack/gin"
	"github.com/go-minstack/sqlite"
	"todo-api/internal/todos"
	todo_entities "todo-api/internal/todos/entities"
	"gorm.io/gorm"
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
