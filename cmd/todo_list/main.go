package main

import (
	"fmt"

	"github.com/KorolevIvanMi/TODO_list_Go/internal/config"
)

func main() {
	//init config
	cfg := config.MustLoad()
	fmt.Println(cfg)
	//TODO : init logger

	//TODO : init storage

	//TODO : init router

	//TODO : run server
}
