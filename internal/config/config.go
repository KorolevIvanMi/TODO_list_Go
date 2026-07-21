package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env          string `yaml:"env" env-required:"true"`
	Storage_path string `yaml:"storage_path" env-required:"true"`
	HTTPServer   `yaml:"http_server"`
}

type HTTPServer struct {
	Adress      string        `yaml:"adress"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not valid")
	}

	//check if file exist
	_, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Error with config file: %s", err)
	}
	var cfg Config
	err = cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Error with reading config: %s", err)
	}

	return &cfg
}
