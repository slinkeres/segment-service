package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL string
	Port string
}

func Load() (*Config, error){
	_ = godotenv.Load()
	env := os.Getenv("ENV")
	var dbURL string
	switch env {
		case "localhost":
			dbURL = os.Getenv("POSTGRES_STRING_LOCALHOST")
		case "docker":
			dbURL = os.Getenv("POSTGRES_STRING_DOCKER")
		default:
			dbURL = os.Getenv("POSTGRES_STRING_DOCKER")
	}
	port := os.Getenv("PORT")
	if port == ""{
		port = "8080"
	}
	return &Config{
		DBURL: dbURL,
		Port: port,
	}, nil
}