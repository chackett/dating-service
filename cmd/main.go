package main

import (
	"github.com/caarlos0/env"
	"github.com/chackett/dating-service/datingservice"
	"github.com/chackett/dating-service/httpserver"
	"github.com/chackett/dating-service/repository"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		logger.Error("Error parsing env config variables: %w", err)
		os.Exit(1)
	}

	repo, err := repository.New(cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)
	if err != nil {
		logger.Error("unable to instantiate repository: %w", err)
		os.Exit(1)
	}

	ds, err := datingservice.New(repo)
	if err != nil {
		logger.Error("unable to instantiate dating service: %w", err)
		os.Exit(1)
	}

	server, err := httpserver.New(cfg.ServicePort, ds)
	if err != nil {
		logger.Error("unable to instantiate http server: %w", err)
		os.Exit(1)
	}

	err = server.Serve()
	if err != nil {
		logger.Error("start webserver", err)
	}
}
