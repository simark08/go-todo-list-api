package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
	"todo-list-echo/auth"
	"todo-list-echo/data"

	"github.com/labstack/echo/v4"
)

type PostRegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PostLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) HandlePostRegister(c echo.Context) error {
	var req PostRegisterRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{
			"message": "Invalid request body",
		})
	}

	exists, err := data.UserExistsByEmail(h.DB, req.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Something went wrong. Please try again",
		})
	}

	if exists {
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "This email address is already taken",
		})
	}

	tx, err := h.DB.Begin()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Something went wrong. Please try again",
		})
	}
	defer tx.Rollback()

	u, err := data.CreateUser(tx, data.CreateUserParams{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Something went wrong. Please try again",
		})
	}

	token, err := auth.GenerateToken(auth.TokenLength)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Something went wrong. Please try again",
		})

	}

	pat, err := data.CreatePersonalAccessToken(tx, data.CreatePersonalAccessTokenParams{
		UserID: u.ID,
		Token:  token,
		ExpiresAt: sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		},
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Something went wrong. Please try again",
		})
	}

	if err := tx.Commit(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Something went wrong. Please try again",
		})
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"user":  u,
		"token": pat.Token,
	})
}

func (h *Handler) HandlePostLogin(c echo.Context) error {
	var req PostLoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{
			"message": "Invalid request body",
		})
	}

	u, err := data.GetUserByEmail(h.DB, req.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Something went wrong. Please try again",
		})
	}

	if u == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid email or password",
		})
	}

	if u.Password != req.Password {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid email or password",
		})
	}

	token, err := auth.GenerateToken(auth.TokenLength)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Something went wrong. Please try again",
		})
	}

	tx, err := h.DB.Begin()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Something went wrong. Please try again",
		})
	}
	defer tx.Rollback()

	pat, err := data.CreatePersonalAccessToken(tx, data.CreatePersonalAccessTokenParams{
		UserID: u.ID,
		Token:  token,
		ExpiresAt: sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		},
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Something went wrong. Please try again",
		})
	}

	if err := tx.Commit(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Something went wrong. Please try again",
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"user":  u,
		"token": pat.Token,
	})
}
