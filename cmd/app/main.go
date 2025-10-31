package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/olenka-91/DocsServer/internal/config"
	"github.com/olenka-91/DocsServer/internal/handler"
	"github.com/olenka-91/DocsServer/internal/repository"
	"github.com/olenka-91/DocsServer/internal/service"
	"github.com/olenka-91/DocsServer/internal/storage"
	"github.com/olenka-91/DocsServer/pkg/httpserver"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Info("Loading environment variables...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	log.Info("Creating DB connection...")
	db, err := repository.NewPostgresDB(*cfg)
	if err != nil {
		log.WithField("err:", err.Error()).Error("Couldn't create DB connection!")
		return
	}
	log.Debug("DB connected successfully")

	log.Info("Creating repositories...")
	repos := repository.NewRepository(db)
	log.Debug("Repositories created successfully")

	log.Info("Creating FileStorage...")
	fs := storage.NewFileStorage(cfg.StorageAddr)
	log.Debug("FileStorage created successfully")

	log.Info("Creating services...")
	serv := service.NewService(repos, fs)
	log.Debug("Services created successfully")

	log.Info("Creating handlers...")
	handl := handler.NewHandler(serv)
	log.Debug("Handlers created successfully")

	log.Info("Creating server...")
	server := new(httpserver.Server)
	log.Debug("Server created successfully")

	go func() {
		log.Info("Starting the HTTP server...")
		if err := server.Run(cfg.HTTPPort, handl.InitRoutes(), db); err != nil {
			if err != http.ErrServerClosed {
				log.Fatalf("error occured while running http server: %s", err.Error())
			}
			log.Info("Server stopped running")
		}
	}()

	log.Info("DocsServer-app Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Info("Shutting down the server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		log.Errorf("error occured on db connection close: %s", err.Error())
	}
	log.Info("Server gracefully stopped")
}
