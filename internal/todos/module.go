package todos

import (
	"github.com/go-minstack/core"
	todo_repos "todo-api/internal/todos/repositories"
)

func Register(app *core.App) {
	app.Provide(todo_repos.NewTodoRepository)
	app.Provide(NewTodoService)
	app.Provide(NewTodoController)
	app.Invoke(RegisterRoutes)
}
