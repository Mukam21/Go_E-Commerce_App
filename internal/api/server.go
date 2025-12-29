package api

import (
	"log"

	"github.com/Mukam21/Go_E-Commerce_App/config"
	"github.com/Mukam21/Go_E-Commerce_App/internal/api/rest"
	"github.com/Mukam21/Go_E-Commerce_App/internal/api/rest/handlers"
	"github.com/Mukam21/Go_E-Commerce_App/internal/domain"
	"github.com/Mukam21/Go_E-Commerce_App/internal/helper"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartServar(config config.AppConfig) {
	app := fiber.New()

	db, err := gorm.Open(postgres.Open(config.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("database connection error %v\n", err)
	}

	log.Printf("database connected")

	// run migrations

	db.AutoMigrate(&domain.User{})

	auth := helper.SetupAuth(config.AppSecret)

	rh := &rest.RestHandler{
		App:    app,
		DB:     db,
		Auth:   auth,
		Config: config,
	}

	setupRoutes(rh)

	app.Listen(config.ServerPort)
}

func setupRoutes(rh *rest.RestHandler) {
	// user handler
	handlers.SetupUserRoutes(rh)
	// transactions
	// catalog
}
