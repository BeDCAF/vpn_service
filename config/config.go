package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type ServiceConfig struct {
	Host  string
	Token string
	Addr  string
	Port  string
}

func Load() (*ServiceConfig, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	host := os.Getenv("HOST")
	token := os.Getenv("TOKEN")
	addr := os.Getenv("ADDR")
	port := os.Getenv("PORT")
	if host == "" || token == "" {
		return nil, errors.New("One of the values is not set, please, check .env")
	}

	return &ServiceConfig{
		Host:  host,
		Token: token,
		Addr:  addr,
		Port:  port,
	}, nil
}
