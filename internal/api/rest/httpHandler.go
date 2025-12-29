package rest

import (
	"github.com/Mukam21/Go_E-Commerce_App/config"
	"github.com/Mukam21/Go_E-Commerce_App/internal/helper"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RestHandler struct {
	App    *fiber.App
	DB     *gorm.DB
	Auth   helper.Auth
	Config config.AppConfig
}
