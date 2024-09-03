package createtaskhandler

import (
	"log/slog"
	"net/http"
	"to-do-list/internal/lib/logger/sl"
	"to-do-list/models/task"

	tr "to-do-list/internal/storage"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func New(log *slog.Logger, taskRepository *tr.TaskRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		newTask := &task.Task{}
		// Декодирование тела запроса
		err := render.DecodeJSON(r.Body, newTask)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			http.Error(w, "failed to decode request", http.StatusBadRequest) // Код 400 Bad Request
			return
		}
		log.Info("request body decoded", slog.Any("newTask", newTask))

		// Сохранение заметки
		savedTask, err := taskRepository.CreateTask(*newTask)
		if err != nil {
			log.Error("failed to save task", sl.Err(err))
			http.Error(w, "failed to save task", http.StatusInternalServerError) // Код 500 Internal Server Error
			return
		}

		// Успешное завершение
		log.Info("task added successfully", slog.Any("task", savedTask))
		w.WriteHeader(http.StatusCreated) // Код 201 Created
		render.JSON(w, r, savedTask)
	}
}
