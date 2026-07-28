package taskhandler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/KorolevIvanMi/TODO_list_Go/adapters/http/dto"
	createtask "github.com/KorolevIvanMi/TODO_list_Go/internal/usecase/taskUsecase/createTask"
	getalltasks "github.com/KorolevIvanMi/TODO_list_Go/internal/usecase/taskUsecase/getAllTasks"
	"github.com/go-chi/render"
)

type TaskHandler struct {
	createUC createtask.UseCase
	getAllUC getalltasks.UseCase
}

func New(
	createUC *createtask.UseCase,
	getAllUC *getalltasks.UseCase,
) *TaskHandler {
	th := TaskHandler{createUC: *createUC, getAllUC: *getAllUC}
	return &th
}

func (handler *TaskHandler) CreateTaskHandler(log slog.Logger) http.HandlerFunc {
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

		task, err := handler.createUC.CreateTask(ctx, req.NAME, req.DESCRIPTION, req.DEADLINE)
		if err != nil {
			log.Error("Failed to create task", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, dto.CreateTaskResponse{STATUS: "ERROR", ERROR: err.Error()})
			return
		}
		log.Info("Creating complited")
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, dto.CreateTaskResponse{STATUS: "OK", NAME: task.Name, ID: task.ID})

	}
}

func (handler *TaskHandler) GetAllTasks(log slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		const op = "internal.delivery.http.handler.task_handler.GetAllTasks"
		log.With(slog.String("op", op))

		tasks, err := handler.getAllUC.GetAllTasks(ctx)
		if err != nil {
			log.Error("Fail to get tasks", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, dto.GetAllTasksResponse{STATUS: "ERROR"})
			return
		}
		if tasks == nil {
			log.Error("Tasks is nil")
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, dto.GetAllTasksResponse{STATUS: "ERROR"})
			return
		}
		var result_task dto.GetAllTasksRespModel

		var response dto.GetAllTasksResponse
		for _, value := range *tasks {
			log.Info("Current task id", slog.String("id", fmt.Sprintf("%d", value.Id)))
			result_task.ID = value.Id
			result_task.DEADLINE = value.Deadline
			result_task.DESCRIPTION = value.Description
			result_task.NAME = value.Name
			response.TASKS = append(response.TASKS, result_task)
		}
		response.AMOUNT = uint64(len(*tasks))
		response.STATUS = "OK"
		log.Info("Parsing completed")
		w.WriteHeader(http.StatusOK)

		render.JSON(w, r, response)
	}
}
