package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/SKjustSK/alru-url-shortener/backend/internal/database"
	"github.com/SKjustSK/alru-url-shortener/backend/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v5"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateUserResponse strictly defines the JSON payload sent back to the client
type CreateUserResponse struct {
	Message string      `json:"message"`
	User    models.User `json:"user"`
}

func CreateUser(c *echo.Context) error {
	req := new(CreateUserRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	// 1. Basic Validation
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.Email == "" || len(req.Password) < 8 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Valid email and password (min 8 characters) are required",
		})
	}

	// 2. Hash the Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to secure password"})
	}

	// 3. Setup User Record
	newUser := models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	// 4. Save to Database
	if err := database.DB.Create(&newUser).Error; err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Email is already registered"})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	// 5. Return the strongly-typed response
	return c.JSON(http.StatusCreated, CreateUserResponse{
		Message: "User created successfully",
		User:    newUser,
	})
}
