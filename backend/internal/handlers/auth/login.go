package auth

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/SKjustSK/alru-url-shortener/backend/internal/database"
	"github.com/SKjustSK/alru-url-shortener/backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

func CreateSession(c *echo.Context) error {
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// 1. Find the User
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
	}

	// 2. Verify the Password
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
	}

	// 3. Create the JWT Payload (Claims)
	claims := &models.JWTCustomClaims{
		UserID: user.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			// Token expires in 24 hours
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	// 4. Create and Sign the Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Get secret from ENV, fallback to a hardcoded string for local dev
	secret := os.Getenv("JWT_SECRET")

	// Sign the token mathematically
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	// 5. Return the Token
	return c.JSON(http.StatusOK, LoginResponse{
		Message: "Login successful",
		Token:   signedToken,
	})
}
