package taskhandler

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/KorolevIvanMi/TODO_list_Go/internal/delivery/http/dto"
	"github.com/KorolevIvanMi/TODO_list_Go/internal/usecase"
	"github.com/go-chi/render"
)

type TaskHandler struct {
	Uc usecase.UseCase
}

func New(uc usecase.UseCase) *TaskHandler {
	th := TaskHandler{Uc: uc}
	return &th
}

func (th *TaskHandler) CreateTaskHandler(log slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		const op = "internal.delivery.http.handler.task_handler.CreateTaskHandler"
		log.With(slog.String("op", op))
		var req dto.CreateTaskRequst

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {

			log.Error("Failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.CreateTaskResponse{STATUS: "ERROR", ERROR: err.Error()})

			return
		}
		log.Info("request body decoded", slog.Any("request", req))

		if req.NAME == "" {
			log.Error("Bad task name", slog.String("name", req.NAME))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.CreateTaskResponse{STATUS: "ERROR", ERROR: "Bad task name"})
			return
		}
		if req.DESCRIPTION == "" {
			log.Error("Bad description", slog.String("desc", req.DESCRIPTION))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.CreateTaskResponse{STATUS: "ERROR", ERROR: "Bad description"})
			return
		}
		log.Info("handler validation is succed")

		task, err := th.Uc.CreateTask(ctx, req.NAME, req.DESCRIPTION, req.DEADLINE)
		if err != nil {
			log.Error("Failed to decode request body", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, dto.CreateTaskResponse{STATUS: "ERROR", ERROR: err.Error()})
			return
		}
		log.Info("Creating complited")
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, dto.CreateTaskResponse{STATUS: "OK", NAME: task.Name, ID: task.ID})

	}
}
