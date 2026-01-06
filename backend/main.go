package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"

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

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisDB := 0
	if redisDBStr := os.Getenv("REDIS_DB"); redisDBStr != "" {
		parsed, err := strconv.Atoi(redisDBStr)
		if err != nil {
			log.Fatalf("invalid REDIS_DB: %v", err)
		}
		redisDB = parsed
	}

	sessionTTL := 24 * time.Hour
	if ttlStr := os.Getenv("SESSION_TTL_SECONDS"); ttlStr != "" {
		seconds, err := strconv.Atoi(ttlStr)
		if err != nil || seconds <= 0 {
			log.Fatalf("invalid SESSION_TTL_SECONDS: %s", ttlStr)
		}
		sessionTTL = time.Duration(seconds) * time.Second
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

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       redisDB,
	})
	defer redisClient.Close()

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	todoUsecase := uc.NewTodoUsecase()
	userUsecase := uc.NewUserUsecase(db)
	sessionUsecase := uc.NewSessionUsecase(redisClient, sessionTTL)

	e.POST("/register", func(c echo.Context) error {
		var req registerRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
		}

		userID, err := userUsecase.Register(req.Name, req.Email, req.Password)
		if err != nil {
			if err == uc.ErrEmailAlreadyExists {
				return echo.NewHTTPError(http.StatusBadRequest, "email already exists")
			}
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}

		sessionID, err := sessionUsecase.Create(c.Request().Context(), userID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}

		cookie := &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Secure:   false,
			Expires:  time.Now().Add(sessionTTL),
		}
		c.SetCookie(cookie)

		return c.NoContent(http.StatusNoContent)
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
