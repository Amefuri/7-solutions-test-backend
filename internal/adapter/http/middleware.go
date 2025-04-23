package http

import (
	"fmt"
	"time"

	"7-solutions-test-backend/internal/auth"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		fmt.Printf("🛠️  %s %s in %v\n", c.Request().Method, c.Path(), time.Since(start))
		return err
	}
}

func AuthMiddleware(jwtService *auth.JWTService) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(jwtService.GetSecret()),
	})
}
