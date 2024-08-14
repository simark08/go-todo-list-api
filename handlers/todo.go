package handlers

import (
	"net/http"
	"strconv"
	"todo-list-echo/data"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleGetTodos(c echo.Context) error {
	todos, err := data.GetTodos(h.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Something went wrong. Please try again",
		})
	}

	return c.JSON(http.StatusOK, todos)
}

func (h *Handler) HandleGetTodo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}

	todo, err := data.GetTodoByID(h.DB, int64(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Something went wrong. Please try again",
		})
	}
	if todo == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Todo not found",
		})
	}

	return c.JSON(http.StatusOK, todo)
}

type PostTodosRequest struct {
	Title string `json:"title"`
}

func (h *Handler) HandlePostTodos(c echo.Context) error {
	var req PostTodosRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{
			"message": "Invalid request body",
		})
	}

	todo, err := data.CreateTodo((h.DB), data.CreateTodoParams{
		Title: req.Title,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Something went wrong. Please try again",
		})
	}

	return c.JSON(http.StatusCreated, todo)

}

type PutTodosRequest struct {
	Title string `json:"title"`
}

func (h *Handler) HandlePutTodos(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}

	var req PutTodosRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{
			"message": "Invalid request body",
		})
	}

	todo, err := data.GetTodoByID(h.DB, int64(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Something went wrong. Please try again",
		})
	}

	if todo == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Todo not found",
		})
	}

	err = todo.Update(h.DB, data.UpdateTodoParams{
		Title: req.Title,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Something went wrong. Please try again",
		})
	}

	return c.JSON(http.StatusOK, todo)
}

func (h *Handler) HandleDeleteTodos(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}

	todo, err := data.GetTodoByID(h.DB, int64(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Something went wrong. Please try again",
		})
	}
	if todo == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Todo not found",
		})
	}

	err = data.DeleteTodoByID(h.DB, int64(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Something went wrong. Please try again",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Todo deleted",
	})
}
