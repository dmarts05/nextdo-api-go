package server

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
)

type Server struct {
	db     *pgxpool.Pool
	router *echo.Echo
	cfg    config
}

func New() (*Server, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}

	return &Server{
		db:     nil,
		router: echo.New(),
		cfg:    cfg,
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	// Database
	pgxConfig, err := pgxpool.ParseConfig(s.cfg.databaseURL)
	if err != nil {
		return err
	}
	pgxConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		return nil
	}

	db, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping(ctx)
	if err != nil {
		return err
	}

	s.db = db

	// Router
	s.loadValidator()
	s.loadMiddleware()
	s.loadRoutes()

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
