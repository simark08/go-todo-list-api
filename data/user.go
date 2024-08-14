package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"todo-list-echo/types"
)

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserParams struct {
	Name     string
	Email    string
	Password string
}

func CreateUser(tx *sql.Tx, params CreateUserParams) (*User, error) {
	query := "INSERT INTO users (name, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	ts := time.Now().UTC()

	res, err := tx.ExecContext(context.Background(), query, params.Name, params.Email, params.Password, ts, ts)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	u, err := GetUserByID(tx, id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func GetUserByID(dbtx types.DBTX, id int64) (*User, error) {
	var u User
	query := "SELECT *FROM users WHERE id = ?"

	if err := dbtx.QueryRowContext(context.Background(), query, id).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func GetUserByEmail(dbtx types.DBTX, email string) (*User, error) {
	var u User
	query := "SELECT *FROM users WHERE email = ?"

	if err := dbtx.QueryRowContext(context.Background(), query, email).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func UserExistsByEmail(db *sql.DB, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
	SELECT 1
	FROM users
	WHERE email = ?
	)`
	if err := db.QueryRowContext(context.Background(), query, email).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}
