package main

import (
	"log/slog"
	"net/http"
	"os"

	//"time"
	"to-do-list/internal/config"
	"to-do-list/internal/lib/logger/sl"
	"to-do-list/internal/storage/postgresql"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	//"url-shorter/models/task"
	createtaskhandler "to-do-list/internal/http-server/handlers/createTaskHandler"
	// "to-do-list/internal/http-server/handlers/get"
	// "to-do-list/internal/http-server/handlers/getusers"
	tr "to-do-list/internal/storage"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	// logger
	log := setupLogger(cfg.Env)

	log.Info("starting noter", slog.String("env", cfg.Env))

	storage, err := postgresql.New(cfg.StorangePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	taskRepository := tr.New(storage)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	router.Post("/tasks", createtaskhandler.New(log, taskRepository))

	router.Get("/tasks", get.New(log, storage))
	router.Get("/tasks/{id}", get.New(log, storage)) 
	router.Put("/tasks/{id}", get.New(log, storage))
	router.Delete("/tasks/{id}", get.New(log, storage))

	// router.Get("/users", getusers.New(log, storage))
	log.Info("server starting", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger { // зависит от того, где запускается, разные уровни

	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}


func (s *Server) GetTaskHandler(w http.ResponseWriter, r *http.Request) { ... }
func (s *Server) CreateTaskHandler(w http.ResponseWriter, r *http.Request) { ... }
func (s *Server) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) { ... }
func (s *Server) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) { ... }