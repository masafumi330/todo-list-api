package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	uc "masafumi330/todo-list-api/usecase"
)

type todoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "3306"
	}
	if dbUser == "" {
		dbUser = "todo_user"
	}
	if dbPassword == "" {
		dbPassword = "todo_password"
	}
	if dbName == "" {
		dbName = "todo_db"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	todoUsecase := uc.NewTodoUsecase()
	userUsecase := uc.NewUserUsecase(db)

	e.POST("/register", func(c echo.Context) error {
		var req registerRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
		}

		if err := userUsecase.Register(req.Name, req.Email, req.Password); err != nil {
			if err == uc.ErrEmailAlreadyExists {
				return echo.NewHTTPError(http.StatusBadRequest, "email already exists")
			}
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}

		return c.NoContent(http.StatusCreated)
	})

	e.POST("/todos", func(c echo.Context) error {
		var req todoRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
		}

		created := todoUsecase.Create(req.Title, req.Description)
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

		updated, ok := todoUsecase.Update(id, req.Title, req.Description)
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

		if !todoUsecase.Delete(id) {
			return echo.NewHTTPError(http.StatusNotFound, "todo not found")
		}

		return c.NoContent(http.StatusNoContent)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
