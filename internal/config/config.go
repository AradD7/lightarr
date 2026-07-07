package config

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"

	"github.com/AradD7/lightarr/internal/database"
	"github.com/AradD7/lightarr/internal/logging"
	"github.com/AradD7/lightarr/internal/rules"
	"github.com/AradD7/lightarr/internal/wiz"
	"github.com/joho/godotenv"
)

type Config struct {
	Db           *database.Queries
	Conn         *net.UDPConn
	BulbsMap     map[string]*wiz.Bulb
	Rules        []rules.Rule
	PlexToken    string
	PlexClientId string
	Logger       *slog.Logger
}

func LoadConfig() *Config {
	logger := logging.InitLogger()

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

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to load DB: %v", err)
	}

	conn, _ := net.ListenUDP("udp", &net.UDPAddr{Port: 38899})
	logger.Info(fmt.Sprintf("Opened a UDP connection on %s", conn.LocalAddr().String()))
	defer conn.Close()

	config := Config{
		Db:           database.New(db),
		Conn:         conn,
		PlexToken:    xPlexToken,
		PlexClientId: xClientId,
		Logger:       logger,
	}

	return &config
}
