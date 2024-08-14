package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateTodoParams struct {
	Title string `json:"title"`
}

type UpdateTodoParams struct {
	Title string `json:"title"`
}

func CreateTodo(db *sql.DB, params CreateTodoParams) (*Todo, error) {
	query := `INSERT INTO todos (title, created_at, updated_at) VALUES (?, ?, ?)`
	ts := time.Now().UTC()

	res, err := db.ExecContext(context.Background(), query, params.Title, ts, ts)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	t, err := GetTodoByID(db, id)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func GetTodos(db *sql.DB) ([]Todo, error) {
	var todos []Todo
	query := `SELECT id, title, created_at, updated_at from todos`
	rows, err := db.QueryContext(context.Background(), query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Todo{}, nil
		}

		return nil, err
	}

	for rows.Next() {
		var t Todo

		err := rows.Scan(&t.ID, &t.Title, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning rows: %v", err)
		}

		todos = append(todos, t)

	}

	return todos, nil
}

func GetTodoByID(db *sql.DB, id int64) (*Todo, error) {
	var t Todo
	query := `SELECT * FROM todos WHERE id = ?`

	if err := db.QueryRowContext(context.Background(), query, id).Scan(&t.ID, &t.Title, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return nil, err
	}

	return &t, nil
}

func (t *Todo) Update(db *sql.DB, params UpdateTodoParams) error {
	query := `UPDATE todos SET title = ?, updated_at = ? WHERE id = ?`
	ts := time.Now().UTC()

	_, err := db.ExecContext(context.Background(), query, params.Title, ts, t.ID)
	if err != nil {
		return err
	}

	t.Title = params.Title

	return nil
}

func DeleteTodoByID(db *sql.DB, id int64) error {
	query := `DELETE FROM todos WHERE id = ?`

	_, err := db.ExecContext(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}
