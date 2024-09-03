package gettaskallhandler

import (
	"log/slog"
	"net/http"
	"strconv"
	"to-do-list/internal/lib/logger/sl"

	tr "to-do-list/internal/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func New(log *slog.Logger, taskRepository *tr.TaskRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.gettaskallhandler.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Декодирование тела запроса
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error("invalid task id", sl.Err(err))
			http.Error(w, "Invalid task ID", http.StatusBadRequest) // Код 400 Bad Request
			return
		}
		log.Info("request id decoded", slog.Any("newTask", id))

		// Сохранение заметки
		task, err := taskRepository.GetTaskByID(id)
		if err != nil {
			log.Error("failed to get task", sl.Err(err))
			http.Error(w, "Task not found", http.StatusNotFound) // Код 404 Not Found
			return
		}

		// Успешное завершение
		log.Info("task retrieved successfully", slog.Any("task", task))
		render.JSON(w, r, task)
	}
}
