package handlers

import (
	"time"

	"github.com/labstack/echo/v5"
)

type ShortenRequest struct {
	LongURL   string
	ShortCode string
	ExpiresIn time.Duration
}

type ShortenResponse struct {
	ShortURL  string
	LongURL   string
	ShortCode string
}

func CreateLink(c *echo.Context) {

}

func GetLinks(c *echo.Context) {

}
