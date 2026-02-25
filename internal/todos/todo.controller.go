package todos

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-minstack/web"
	"todo-api/internal/todos/dto"
)

type TodoController struct {
	service *TodoService
}

func NewTodoController(service *TodoService) *TodoController {
	return &TodoController{service: service}
}

func (c *TodoController) list(ctx *gin.Context) {
	todos, err := c.service.List()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, web.NewErrorDto(err))
		return
	}
	ctx.JSON(http.StatusOK, todos)
}

func (c *TodoController) create(ctx *gin.Context) {
	var input dto.CreateTodoDto
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, web.NewErrorDto(err))
		return
	}
	todo, err := c.service.Create(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, web.NewErrorDto(err))
		return
	}
	ctx.JSON(http.StatusCreated, todo)
}

func (c *TodoController) get(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, web.NewErrorDto(err))
		return
	}
	todo, err := c.service.Get(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, web.NewErrorDto(err))
		return
	}
	ctx.JSON(http.StatusOK, todo)
}

func (c *TodoController) update(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, web.NewErrorDto(err))
		return
	}
	var input dto.UpdateTodoDto
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, web.NewErrorDto(err))
		return
	}
	todo, err := c.service.Update(id, input)
	if err != nil {
		ctx.JSON(http.StatusNotFound, web.NewErrorDto(err))
		return
	}
	ctx.JSON(http.StatusOK, todo)
}

func (c *TodoController) delete(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, web.NewErrorDto(err))
		return
	}
	if err := c.service.Delete(id); err != nil {
		ctx.JSON(http.StatusNotFound, web.NewErrorDto(err))
		return
	}
	ctx.Status(http.StatusNoContent)
}

func parseID(ctx *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		return 0, errors.New("invalid id")
	}
	return uint(id), nil
}
