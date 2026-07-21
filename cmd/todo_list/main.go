package main

import (
	"log/slog"
	"os"

	"github.com/KorolevIvanMi/TODO_list_Go/internal/config"
	"github.com/KorolevIvanMi/TODO_list_Go/internal/delivery/sqlite"
	"github.com/KorolevIvanMi/TODO_list_Go/internal/logger"
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
	_ = storage
	//TODO : init router

	//TODO : run server
}
