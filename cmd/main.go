package main

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/chackett/dating-service/datingservice"
	"github.com/chackett/dating-service/httpserver"
	"github.com/chackett/dating-service/repository"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	logger *slog.Logger
	cfg    *Config
	server *httpserver.HTTPServer
}

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

	//app := App{
	//	logger: logger,
	//	cfg:    cfg,
	//	server: server,
	//}

	//go app.listenForSigKill()

	err = server.Serve()
	if err != nil {
		logger.Error("start webserver", err)
	}
}

func cleanup() {

}

func (app *App) listenForSigKill() {
	fmt.Println("listen for kill")
	// Set up signal handling.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	done := make(chan bool, 1)
	go func() {
		sig := <-signals
		fmt.Println("")
		fmt.Println("Disconnection requested via Ctrl+C", sig)
		done <- true
	}()

	app.logger.Debug("SIGKILL detected")
	<-done

	err := app.server.Close()
	if err != nil {
		app.logger.Error("shutdown webserver", err)
		os.Exit(1)
	}
	app.logger.Info("did close webserver")

	cleanup()

	os.Exit(0)
}
