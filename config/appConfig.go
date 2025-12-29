package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	ServerPort   string
	Dsn          string
	AppSecret    string
	Env          string
	SMSRuApiKey  string
	EnableSMSDev bool // флаг для реальной отправки в dev
	TestPhone    string
}

func SetupEnv() (cfg AppConfig, err error) {

	if os.Getenv("APP_ENV") == "dev" {
		godotenv.Load()
	}

	httpPort := os.Getenv("HTTP_PORT")

	if len(httpPort) < 1 {
		return AppConfig{}, errors.New("env variables not found")
	}

	Dsn := os.Getenv("Dsn")
	if len(Dsn) < 1 {
		return AppConfig{}, errors.New("env variables not found")
	}

	appSecret := os.Getenv("APP_SECRET")
	if len(appSecret) < 1 {
		return AppConfig{}, errors.New("app secret not found")
	}

	return AppConfig{
		ServerPort:  httpPort,
		Dsn:         Dsn,
		AppSecret:   appSecret,
		Env:         os.Getenv("ENV"),
		SMSRuApiKey: os.Getenv("SMS_RU_API_KEY"),
	}, nil
}
