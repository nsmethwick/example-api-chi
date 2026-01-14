package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type API struct {
	Address string
	Host    string
	Port    string
}

type DB struct {
	MYSQL_DATABASE      string
	MYSQL_PASSWORD      string
	MYSQL_USER          string
	MYSQL_ROOT_PASSWORD string
	DSN                 string
}

func SetUpAPI() API {
	api := API{
		Host:    os.Getenv("API_HOST"),
		Port:    os.Getenv("API_PORT"),
		Address: os.Getenv("API_ADDRESS"),
	}
	return api
}

func SetUpDB() DB {
	db := DB{
		MYSQL_DATABASE:      os.Getenv("MYSQL_DATABASE"),
		MYSQL_PASSWORD:      os.Getenv("MYSQL_PASSWORD"),
		MYSQL_USER:          os.Getenv("MYSQL_USER"),
		MYSQL_ROOT_PASSWORD: os.Getenv("MYSQL_ROOT_PASSWORD"),
		DSN:                 os.Getenv("MYSQL_USER") + ":" + os.Getenv("MYSQL_PASSWORD") + "@tcp(localhost:3306)/" + os.Getenv("MYSQL_DATABASE") + "?parseTime=true",
	}
	return db
}

func LoadEnvironmentVariables(env string) {
	// This will load .env files based on args
	// 1: local == .env.local
	// 2: staging == .env.staging
	// 3: prod == .env

	// If the files are not found, we will default to using .env
	// If .env is not found, we just proceed without loading a file

	candidates := []string{}
	switch env {
	case "local", "":
		// default to local during development when env is empty
		candidates = []string{".env.local", ".env"}
	case "staging":
		candidates = []string{".env.staging", ".env"}
	case "prod":
		candidates = []string{".env"}
	default:
		candidates = []string{".env"}
	}

	for _, envFile := range candidates {
		if doesFileExist(envFile) {
			slog.Info("Loading Environment Variables from " + envFile)
			if err := godotenv.Load(envFile); err != nil {
				slog.Error("failed to load env file", "file", envFile, "err", err)
			}
			return
		}
	}

	slog.Info("No env file found; proceeding without loading from file")

}

func doesFileExist(fileName string) bool {
	_, error := os.Stat(fileName)

	// check if error is "file not exists"
	if os.IsNotExist(error) {
		return false
	} else {
		return true
	}
}
