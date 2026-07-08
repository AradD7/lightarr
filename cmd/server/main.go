package main

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"

	"github.com/AradD7/lightarr/frontend"
	"github.com/AradD7/lightarr/internal/config"
	"github.com/AradD7/lightarr/internal/database"
	"github.com/AradD7/lightarr/internal/handlers"
	"github.com/AradD7/lightarr/internal/logging"
	"github.com/AradD7/lightarr/internal/rules"
	"github.com/AradD7/lightarr/internal/wiz"
	"github.com/AradD7/lightarr/sql/schema"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

func main() {
	app := handlers.App{
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

	app.Conn = conn

	goose.SetBaseFS(schema.EmbedMigrations)
	if err := goose.SetDialect("sqlite"); err != nil {
		log.Fatalf("Failed to load migrations: %v", err)
	}

	if err := goose.Up(db, "."); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	app.Db = database.New(db)

	frontend, err := fs.Sub(frontend.EmbedStatics, "dist")
	if err != nil {
		log.Fatalf("Failed to load frontend: %v", err.Error())
	}

	app.BulbsMap = wiz.LoadBulbs(conn, app.Db, app.Logger)
	app.Rules, err = rules.LoadRulesFromDb(app.Db, app.Logger)
	if err != nil {
		app.Logger.Info("Failed to load rules", err.Error())
	}

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.FS(frontend)))
	mux.HandleFunc("GET /api/bulbs", app.HandlerGetBulbs)
	mux.HandleFunc("POST /api/bulbs/name", app.HandlerUpdateBulbName)
	mux.HandleFunc("POST /api/bulbs/flash", app.HandlerFlashBulb)
	mux.HandleFunc("GET /api/bulbs/refresh", app.HandlerRefreshBulbs)
	mux.HandleFunc("POST /api/bulbs/type", app.HandlerUpdateBulbType)

	mux.HandleFunc("GET /api/accounts", app.HandlerGetAllAccounts)
	mux.HandleFunc("GET /api/devices", app.HandlerGetAllDevices)
	mux.HandleFunc("POST /api/accounts", app.HandlerAddAccount)
	mux.HandleFunc("POST /api/devices", app.HandlerAddDevice)
	mux.HandleFunc("DELETE /api/accounts/{accountId}", app.HandlerDeleteAccount)
	mux.HandleFunc("DELETE /api/devices/{deviceId}", app.HandlerDeleteDevice)

	mux.HandleFunc("GET /api/rules", app.HandlerGetAllRules)
	mux.HandleFunc("POST /api/rules", app.HandlerAddRule)
	mux.HandleFunc("DELETE /api/rules/{ruleId}", app.HandlerDeleteRule)
	mux.HandleFunc("POST /api/rules/name", app.HandlerUpdateRuleName)

	mux.HandleFunc("POST /plexhook", app.HandlerPlexWebhook)

	mux.HandleFunc("GET /api/plex/accounts", app.HandlerPlexAllAccounts)
	mux.HandleFunc("GET /api/plex/devices", app.HandlerPlexAllDevices)

	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + app.Cfg.Port,
	}

	app.Logger.Info(fmt.Sprintf("Lightarr available on port: %s", app.Cfg.Port))
	srv.ListenAndServe()
}
