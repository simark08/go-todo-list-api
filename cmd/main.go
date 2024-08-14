package main

import (
	"log"
	"os"
	"todo-list-echo/db"
	"todo-list-echo/handlers"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	db := db.New(&db.Config{
		Engine:   os.Getenv("DB_ENGINE"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_DATABASE"),
	})
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	h := handlers.New(db)
	e.GET("/todos", h.HandleGetTodos)
	e.GET("/todos/:id", h.HandleGetTodo)
	e.POST("/todos", h.HandlePostTodos)
	e.PUT("/todos/:id", h.HandlePutTodos)
	e.DELETE("/todos/:id", h.HandleDeleteTodos)

	e.POST("/auth/register", h.HandlePostRegister)
	e.POST("/auth/login", h.HandlePostLogin)

	e.Start(":8000")
}
