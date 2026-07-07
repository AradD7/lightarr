package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/AradD7/lightarr/internal/config"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

//go:embed sql/schema/*.sql
var embedMigrations embed.FS

//go:embed frontend/dist
var embedStatics embed.FS

func main() {
	config := config.LoadConfig()

	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("sqlite"); err != nil {
		log.Fatalf("Failed to load migrations: %v", err)
	}

	if err := goose.Up(config.db, "sql/schema"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	frontend, err := fs.Sub(embedStatics, "frontend/dist")
	if err != nil {
		log.Fatalf("Failed to load frontend: %v", err.Error())
	}

	config.LoadBulbs(conn)
	err = config.loadRules()
	if err != nil {
		logger.Info("Failed to load rules", err.Error())
	}

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.FS(frontend)))
	mux.HandleFunc("GET /api/bulbs", config.handlerGetBulbs)
	mux.HandleFunc("POST /api/bulbs/name", config.handlerUpdateBulbName)
	mux.HandleFunc("POST /api/bulbs/flash", config.handlerFlashBulb)
	mux.HandleFunc("GET /api/bulbs/refresh", config.handlerRefreshBulbs)
	mux.HandleFunc("POST /api/bulbs/type", config.handlerUpdateBulbType)

	mux.HandleFunc("GET /api/accounts", config.handlerGetAllAccounts)
	mux.HandleFunc("GET /api/devices", config.handlerGetAllDevices)
	mux.HandleFunc("POST /api/accounts", config.handlerAddAccount)
	mux.HandleFunc("POST /api/devices", config.handlerAddDevice)
	mux.HandleFunc("DELETE /api/accounts/{accountId}", config.handlerDeleteAccount)
	mux.HandleFunc("DELETE /api/devices/{deviceId}", config.handlerDeleteDevice)

	mux.HandleFunc("GET /api/rules", config.handlerGetAllRules)
	mux.HandleFunc("POST /api/rules", config.handlerAddRule)
	mux.HandleFunc("DELETE /api/rules/{ruleId}", config.handlerDeleteRule)
	mux.HandleFunc("POST /api/rules/name", config.handlerUpdateRuleName)

	mux.HandleFunc("POST /plexhook", config.handlerPlexWebhook)

	mux.HandleFunc("GET /api/plex/accounts", config.handlerPlexAllAccounts)
	mux.HandleFunc("GET /api/plex/devices", config.handlerPlexAllDevices)

	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	logger.Info(fmt.Sprintf("Lightarr available on port: %s", port))
	srv.ListenAndServe()
}
