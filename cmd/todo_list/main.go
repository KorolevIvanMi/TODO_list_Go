package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/KorolevIvanMi/TODO_list_Go/internal/config"
	taskhandler "github.com/KorolevIvanMi/TODO_list_Go/internal/delivery/http/handler/task_handler"
	"github.com/KorolevIvanMi/TODO_list_Go/internal/delivery/sqlite"
	"github.com/KorolevIvanMi/TODO_list_Go/internal/logger"
	"github.com/KorolevIvanMi/TODO_list_Go/internal/usecase"
	"github.com/go-chi/chi/v5"
)

func main() {
	//init config
	cfg := config.MustLoad()

	//init logger
	log := logger.SetupLogger(cfg.Env)
	log.Info("Config is ready")
	log.Info("Logger is ready")
	log.Debug("Config info: ", slog.String("env", cfg.Env), slog.String("adress", cfg.Adress), slog.String("starage path", cfg.Storage_path))

	//init storage
	storage, err := sqlite.New(cfg.Storage_path)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	log.Info("Data base ready")
	var repo sqlite.TaskRepository = storage

	taskUseCase := usecase.New(repo)

	log.Info("UseCase ready")

	taskHandler := taskhandler.New(*taskUseCase)
	log.Info("Task Handler ready")
	//TODO : init router
	router := chi.NewRouter()
	router.Post("/task", taskHandler.CreateTaskHandler(*log))
	log.Info("Router init")
	//TODO : run server
	srv := &http.Server{
		Addr:         cfg.Adress,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	log.Info("Server struct ready")
	log.Info("starting todo list", slog.String("adress", cfg.Adress))

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
		os.Exit(1)
	}
	// log.Error("server stopped")s
}
