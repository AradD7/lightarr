package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net"

	"github.com/AradD7/lightarr/internal/config"
	"github.com/AradD7/lightarr/internal/database"
	"github.com/AradD7/lightarr/internal/logging"
	"github.com/AradD7/lightarr/internal/rules"
	"github.com/AradD7/lightarr/internal/wiz"
)

type App struct {
	Cfg      config.Config
	Db       *database.Queries
	Conn     *net.UDPConn
	BulbsMap map[string]*wiz.Bulb
	Rules    []rules.Rule
	Logger   *slog.Logger
}

func LoadApp() *App {
	app := App{
		Cfg:    config.LoadConfig(),
		Logger: logging.InitLogger(),
	}

	db, err := sql.Open("sqlite", app.Cfg.DbPath)
	if err != nil {
		log.Fatalf("Failed to load DB: %v", err)
	}

	conn, _ := net.ListenUDP("udp", &net.UDPAddr{Port: 38899})
	app.Logger.Info(fmt.Sprintf("Opened a UDP connection on %s", conn.LocalAddr().String()))
	defer conn.Close()

	app.Db = database.New(db)
	app.Conn = conn

	return &app
}
