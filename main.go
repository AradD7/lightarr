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
	db 			*database.Queries
	conn 		*net.UDPConn
	bulbsMap	map[string]*wiz.Bulb
	rules 		[]Rule
}

//go:embed sql/schema/*.sql
var embedMigrations embed.FS

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	dbPath := os.Getenv("DB_PATH")

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
		db: 		database.New(db),
		conn: 		conn,
	}
	config.LoadBulbs(conn)
	config.loadRules()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/bulbs", config.handlerGetBulbs)

	srv := &http.Server {
		Handler: mux,
		Addr: 	 ":" + port,
	}

	fmt.Printf("Api available on port: %s\n", port)
	srv.ListenAndServe()
}

