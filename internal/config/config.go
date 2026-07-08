package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DbPath       string
	PlexToken    string
	PlexClientId string
}

func LoadConfig() Config {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "10100"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/lightarr.db"
	}

	xClientId := os.Getenv("X_PLEX_CLIENT_IDENTIFIER")
	if xClientId == "" {
		xClientId = "lightarr-app"
	}

	xPlexToken := os.Getenv("X_PLEX_TOKEN")

	return Config{
		Port:         port,
		DbPath:       dbPath,
		PlexToken:    xPlexToken,
		PlexClientId: xClientId,
	}
}
