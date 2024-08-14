package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"todo-list-echo/types"
)

type PersonalAccessToken struct {
	ID        int64
	UserID    int64
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt sql.NullTime
}

type CreatePersonalAccessTokenParams struct {
	UserID    int64
	Token     string
	ExpiresAt sql.NullTime
}

func CreatePersonalAccessToken(dbtx types.DBTX, params CreatePersonalAccessTokenParams) (*PersonalAccessToken, error) {
	query := "INSERT INTO personal_access_tokens (user_id, token, created_at, updated_at, expires_at) VALUES (?, ?, ?, ?, ?)"
	ts := time.Now().UTC()

	res, err := dbtx.ExecContext(context.Background(), query, params.UserID, params.Token, ts, ts, params.ExpiresAt)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	pat, err := GetPersonalAccessTokenByID(dbtx, id)
	if err != nil {
		return nil, err
	}

	return pat, nil
}

func GetPersonalAccessTokenByID(dbtx types.DBTX, id int64) (*PersonalAccessToken, error) {
	var pat PersonalAccessToken
	query := "SELECT * FROM personal_access_tokens WHERE id = ?"

	if err := dbtx.QueryRowContext(context.Background(), query, id).Scan(&pat.ID, &pat.UserID, &pat.Token, &pat.CreatedAt, &pat.UpdatedAt, &pat.ExpiresAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &pat, nil
}
