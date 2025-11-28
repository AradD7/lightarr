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

type config struct {
	db 			 *database.Queries
	conn 		 *net.UDPConn
	bulbsMap	 map[string]*wiz.Bulb
	rules 		 []Rule
	plexToken 	 string
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

	conn, _ := net.ListenUDP("udp", &net.UDPAddr{Port: 0})
	fmt.Printf("Opened a UDP connection on %s\n", conn.LocalAddr().String())
	defer conn.Close()


	config := config{
		db: 			database.New(db),
		conn: 			conn,
		plexToken:  	xPlexToken,
		plexClientId: 	xClientId,
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

	mux.HandleFunc("GET /api/accounts", config.handlerGetAllAccounts)
	mux.HandleFunc("GET /api/players", config.handlerGetAllPlayers)
	mux.HandleFunc("POST /api/accounts", config.handlerAddAccounts)
	mux.HandleFunc("POST /api/players", config.handlerAddPlayers)
	mux.HandleFunc("DELETE /api/accounts/{accountId}", config.handlerDeleteAccount)
	mux.HandleFunc("DELETE /api/devices/{playerId}", config.handlerDeletePlayer)

	mux.HandleFunc("GET /api/rules", config.handlerGetAllRules)
	mux.HandleFunc("POST /api/rules", config.handlerAddRule)
	mux.HandleFunc("DELETE /api/rules/{ruleId}", config.handlerDeleteRule)

	mux.HandleFunc("POST /plexhook", config.handlerPlexWebhook)

	mux.HandleFunc("GET /api/plex/accounts", config.handlerPlexAllAccounts)
	mux.HandleFunc("GET /api/plex/players", config.handlerPlexAllPlayers)

	srv := &http.Server {
		Handler: mux,
		Addr: 	 ":" + port,
	}

	fmt.Printf("Api available on port: %s\n", port)
	srv.ListenAndServe()
}

