package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KorolevIvanMi/TODO_list_Go/internal/config"
	taskHandler "github.com/KorolevIvanMi/TODO_list_Go/internal/delivery/http/handler/task_handler"
	"github.com/KorolevIvanMi/TODO_list_Go/internal/delivery/sqlite"
	"github.com/KorolevIvanMi/TODO_list_Go/internal/logger"
	createtask "github.com/KorolevIvanMi/TODO_list_Go/internal/usecase/taskUsecase/createTask"
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
	defer storage.Close()

	log.Info("Data base ready")
	var repo sqlite.TaskRepository = storage

	// taskUseCase := usecase.New(repo)
	log.Info("UseCase ready")

	// taskHandler := taskhandler.New(*taskUseCase)
	log.Info("Task Handler ready")

	// init router
	router := chi.NewRouter()
	router.Post("/task", taskHandler.CreateTaskHandler(*log, createtask.UseCase{Repo: repo}))
	log.Info("Router init")

	// run server
	srv := &http.Server{
		Addr:         cfg.Adress,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	log.Info("Server struct ready")
	log.Info("starting todo list", slog.String("adress", cfg.Adress))

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("failed to start server")
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT)
	<-quit

	log.Info("Shutting down server gracefully...")
	shutdownTimeout := 30 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server shutdown error", slog.String("error", err.Error()))
	}

	log.Info("Server stopped gracefully")

	// log.Error("server stopped")s
}
