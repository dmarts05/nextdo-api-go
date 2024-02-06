package server

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	db     *pgxpool.Pool
	router *echo.Echo
	cfg    config
}

func New() (*Server, error) {
	// Config
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}

	// Router
	r := echo.New()
	r.HideBanner = true
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())

	return &Server{
		db:     nil,
		router: r,
		cfg:    cfg,
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	// Database
	db, err := pgxpool.New(ctx, s.cfg.databaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping(ctx)
	if err != nil {
		return err
	}

	s.db = db

	// Start server in a separate goroutine
	go func() {
		if err := s.router.Start(":" + strconv.Itoa(s.cfg.port)); err != nil {
			s.router.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10)
	defer cancel()
	if err := s.router.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
