package main

import (
	"github.com/KorolevIvanMi/TODO_list_Go/internal/config"
	"github.com/KorolevIvanMi/TODO_list_Go/internal/logger"
)

func main() {
	//init config
	cfg := config.MustLoad()

	//TODO : init logger
	log := logger.SetupLogger(cfg.Env)
	log.Info("Config is ready")
	log.Info("Logger is ready")

	//TODO : init storage

	//TODO : init router

	//TODO : run server
}
