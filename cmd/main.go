package main

import (
	"log/slog"
	"os"
	"time"

	//"time"
	"url-shorter/internal/config"
	"url-shorter/internal/lib/logger/sl"
	"url-shorter/internal/storage/postgresql"
	"url-shorter/models/task"
	//"url-shorter/models/task"
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
	task := &task.Task{
		Title:       "Любой заголовок",
		Description: "Описание задачи",
		DueDate:     time.Now(),
	}

	err = storage.SaveTask(task)

	if err != nil {
		log.Error("failed to SaveTask", sl.Err(err))
		os.Exit(1)
	}

	// router := chi.NewRouter()

	// router.Use(middleware.RequestID)
	// router.Use(middleware.Logger)
	// router.Use(mwLogger.New(log))

	// router.Post("/", save.New(log, storage))
	// router.Get("/", get.New(log, storage))
	// router.Get("/users", getusers.New(log, storage))
	// log.Info("server starting", slog.String("address", cfg.Address))

	// srv := &http.Server{
	// 	Addr:         cfg.Address,
	// 	Handler:      router,
	// 	ReadTimeout:  cfg.HTTPServer.Timeout,
	// 	WriteTimeout: cfg.HTTPServer.Timeout,
	// 	IdleTimeout:  cfg.IdleTimeout,
	// }

	// if err := srv.ListenAndServe(); err != nil {
	// 	log.Error("failed to start server")
	// }

	// log.Error("server stopped")
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
