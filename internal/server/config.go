package server

import (
	"os"
	"strconv"
)

type config struct {
	databaseURL string
	port        int
}

func loadConfig() (config, error) {
	cfg := config{
		databaseURL: "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
		port:        1323,
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL != "" {
		cfg.databaseURL = databaseURL
	}

	portStr := os.Getenv("PORT")
	if portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return cfg, err
		}

		cfg.port = port
	}

	return cfg, nil
}
