package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/AradD7/lightarr/internal/database"
	"github.com/AradD7/lightarr/internal/wiz"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Continue to the next handler
		next.ServeHTTP(w, r)
	})
}

type config struct {
	db           *database.Queries
	conn         *net.UDPConn
	bulbsMap     map[string]*wiz.Bulb
	rules        []Rule
	plexToken    string
	plexClientId string
}

//go:embed sql/schema/*.sql
var embedMigrations embed.FS

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	dbPath := os.Getenv("DB_PATH")
	xPlexToken := os.Getenv("X_PLEX_TOKEN")
	xClientId := os.Getenv("X_PLEX_CLIENT_IDENTIFIER")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to load DB: %v", err)
	}
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("sqlite"); err != nil {
		log.Fatalf("Failed to load migrations: %v", err)
	}

	if err := goose.Up(db, "sql/schema"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	conn, _ := net.ListenUDP("udp", &net.UDPAddr{Port: 38899})
	fmt.Printf("Opened a UDP connection on %s\n", conn.LocalAddr().String())
	defer conn.Close()

	config := config{
		db:           database.New(db),
		conn:         conn,
		plexToken:    xPlexToken,
		plexClientId: xClientId,
	}
	config.LoadBulbs(conn)
	err = config.loadRules()
	if err != nil {
		fmt.Println("Failed to load rules", err.Error())
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/bulbs", config.handlerGetBulbs)
	mux.HandleFunc("POST /api/bulbs/updatename", config.handlerUpdateBulbName)
	mux.HandleFunc("POST /api/bulbs/flash", config.handlerFlashBulb)
	mux.HandleFunc("GET /api/bulbs/refresh", config.handlerRefreshBulbs)

	mux.HandleFunc("GET /api/accounts", config.handlerGetAllAccounts)
	mux.HandleFunc("GET /api/devices", config.handlerGetAllDevices)
	mux.HandleFunc("POST /api/accounts", config.handlerAddAccount)
	mux.HandleFunc("POST /api/devices", config.handlerAddDevice)
	mux.HandleFunc("DELETE /api/accounts/{accountId}", config.handlerDeleteAccount)
	mux.HandleFunc("DELETE /api/devices/{deviceId}", config.handlerDeleteDevice)

	mux.HandleFunc("GET /api/rules", config.handlerGetAllRules)
	mux.HandleFunc("POST /api/rules", config.handlerAddRule)
	mux.HandleFunc("DELETE /api/rules/{ruleId}", config.handlerDeleteRule)

	mux.HandleFunc("POST /plexhook", config.handlerPlexWebhook)

	mux.HandleFunc("GET /api/plex/accounts", config.handlerPlexAllAccounts)
	mux.HandleFunc("GET /api/plex/devices", config.handlerPlexAllDevices)

	srv := &http.Server{
		Handler: corsMiddleware(mux),
		Addr:    ":" + port,
	}

	fmt.Printf("Api available on port: %s\n", port)
	srv.ListenAndServe()
}
