package http

import (
	"7-solutions-test-backend/internal/auth"
	"7-solutions-test-backend/internal/core/user"
	"7-solutions-test-backend/internal/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service    *user.Service
	jwtService *auth.JWTService
}

func NewHandler(service *user.Service, jwt *auth.JWTService) *Handler {
	return &Handler{service, jwt}
}

func (h *Handler) RegisterRoutes(e *echo.Echo) {
	e.POST("api/v1/register", h.Register)
	e.POST("api/v1/login", h.Login)

	r := e.Group("api/v1/users", AuthMiddleware(h.jwtService))
	r.GET("", h.ListUsers)
	r.GET("/:id", h.GetUser)
	r.PUT("/:id", h.UpdateUser)
	r.DELETE("/:id", h.DeleteUser)
}

func (h *Handler) Register(c echo.Context) error {
	req := new(struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	})
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	// Validate required fields
	if req.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Name is required")
	}
	if req.Email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Email is required")
	}
	if req.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Password is required")
	}

	// Validate email format
	if !util.ValidateEmail(req.Email) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid email format")
	}

	// Validate password length
	if len(req.Password) < 8 {
		return echo.NewHTTPError(http.StatusBadRequest, "Password must be at least 8 characters long")
	}

	user, err := h.service.Register(c.Request().Context(), req.Name, req.Email, req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, user)
}

func (h *Handler) Login(c echo.Context) error {
	req := new(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	})
	if err := c.Bind(req); err != nil {
		return err
	}
	user, err := h.service.Authenticate(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	token, _ := h.jwtService.GenerateToken(user.ID)
	return c.JSON(http.StatusOK, echo.Map{"token": token})
}

func (h *Handler) ListUsers(c echo.Context) error {
	users, _ := h.service.List(c.Request().Context())
	return c.JSON(http.StatusOK, users)
}

func (h *Handler) GetUser(c echo.Context) error {
	id := c.Param("id")
	user, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateUser(c echo.Context) error {
	id := c.Param("id")
	req := new(user.User)
	if err := c.Bind(req); err != nil {
		return err
	}
	req.ID = id

	// Validate required fields
	if req.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Name is required")
	}
	if req.Email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Email is required")
	}

	// Validate email format
	if !util.ValidateEmail(req.Email) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid email format")
	}

	if err := h.service.Update(c.Request().Context(), req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
