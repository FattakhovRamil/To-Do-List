package get

import (
	"fmt"
	"log/slog"
	"net/http"
	"url-shorter/internal/lib/logger/sl"
	"url-shorter/models/task"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type User struct {
	name string
	role string
	hash string // пока просто так оставим, заменим потом на реальное хэширование
}

type Request struct {
	Name string `json:"name" validate:"required"`
	Hash string `json:"hash" validate:"required"`
}

type ErrorResponse struct {
	Message string       `json:"message"`
	Errors  []SpellError `json:"errors"`
}

type NotesResponse struct {
	Status  string       `json:"status"`
	Tasks   []*task.Task `json:"tasks,omitempty"`
	Message string       `json:"message,omitempty"`
}

type GetNotesI interface {
	GetNotes(id int) ([]*task.Task, error)
}

func New(log *slog.Logger, noteGetter GetNotesI) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.get.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request

		// Декодирование тела запроса
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			http.Error(w, "failed to decode request", http.StatusBadRequest) // Код 400 Bad Request
			return
		}
		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			// Создание списка сообщений об ошибках
			var errorMessages []string
			for _, fieldError := range validateErr {
				errorMessage := fmt.Sprintf("Field '%s' is required", fieldError.Field())
				errorMessages = append(errorMessages, errorMessage)
			}

			// Логирование ошибки
			log.Error("invalid request", sl.Err(err))

			// Формирование и отправка ответа с ошибками
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, map[string]interface{}{
				"status":  "error",
				"message": "validation failed",
				"errors":  errorMessages,
			})

			return
		}

		// Создание задачи


		// Получение списка задач
		tasks, err := noteGetter.GetNotes(task.UserID)
		if err != nil {
			log.Error("failed to get notes", sl.Err(err))
			http.Error(w, "failed to get notes", http.StatusInternalServerError) // Код 500 Internal Server Error
			return
		}

		// Успешное завершение
		response := NotesResponse{
			Status: "success",
			Tasks:  tasks,
		}

		render.JSON(w, r, response)
		log.Info("Notes sent successfully")
	}
}


