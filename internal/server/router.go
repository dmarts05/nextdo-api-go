package server

import (
	"net/http"

	"github.com/dmarts05/nextdo-api-go/internal/auth"
	"github.com/dmarts05/nextdo-api-go/internal/customer"
	"github.com/labstack/echo/v4"
)

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
	g.POST("/refresh", h.Refresh)
}
