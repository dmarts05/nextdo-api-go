package server

import (
	"net/http"

	"github.com/dmarts05/nextdo-api-go/internal/auth"
	"github.com/dmarts05/nextdo-api-go/internal/customer"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) loadMiddleware() {
	s.router.Use(middleware.Logger())
	s.router.Use(middleware.Recover())
}

func (s *Server) loadValidator() {
	s.router.Validator = &customValidator{validator: validator.New()}
}

func (s *Server) loadRoutes() {
	s.router.GET("/", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	s.loadAuthRoutes()
}

func (s *Server) loadAuthRoutes() {
	h := auth.Handler{
		CustomerRepo: customer.PostgresRepository{
			Db: s.db,
		},
	}

	g := s.router.Group("/auth")
	g.POST("/login", h.Login)
	g.POST("/register", h.Register)

	// TODO: These routes should be protected
	g.POST("/refresh", h.Refresh)
	g.POST("/logout", h.Logout)
}
