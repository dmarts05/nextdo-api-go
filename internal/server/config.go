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
		port:        8080,
	}

	databaseURL, ok := os.LookupEnv("DATABASE_URL")
	if ok {
		cfg.databaseURL = databaseURL
	}

	portStr, ok := os.LookupEnv("PORT")
	if ok {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return cfg, err
		}

		cfg.port = port
	}

	return cfg, nil
}
