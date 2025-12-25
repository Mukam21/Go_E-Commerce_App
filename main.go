package main

import (
	"log"

	"github.com/Mukam21/Go_E-Commerce_App/config"
	"github.com/Mukam21/Go_E-Commerce_App/internal/api"
)

func main() {

	cfg, err := config.SetupEnv()
	if err != nil {
		log.Fatalf("config file is not loaded properly %v\n", err)
	}

	api.StartServar(cfg)
}
