package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	uc "masafumi330/todo-list-api/usecase"
)

type todoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	usecase := uc.NewTodoUsecase()

	e.POST("/todos", func(c echo.Context) error {
		var req todoRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
		}

		created := usecase.Create(req.Title, req.Description)
		return c.JSON(http.StatusCreated, created)
	})

	e.PUT("/todos/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id <= 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
		}

		var req todoRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
		}

		updated, ok := usecase.Update(id, req.Title, req.Description)
		if !ok {
			return echo.NewHTTPError(http.StatusNotFound, "todo not found")
		}

		return c.JSON(http.StatusOK, updated)
	})

	e.DELETE("/todos/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id <= 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
		}

		if !usecase.Delete(id) {
			return echo.NewHTTPError(http.StatusNotFound, "todo not found")
		}

		return c.NoContent(http.StatusNoContent)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
