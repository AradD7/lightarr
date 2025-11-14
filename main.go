package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/AradD7/lightarr/internal/wiz"
)

type config struct {
	conn 		*net.UDPConn
	bulbsMap	map[string]*wiz.Bulb
	rules 		[]Rule
}

const port = "10100"

func main() {
	conn, _ := net.ListenUDP("udp", &net.UDPAddr{Port: 0})
	fmt.Printf("Opened a UDP connection on %s\n", conn.LocalAddr().String())
	defer conn.Close()

	bulbs := wiz.LoadBulbs(conn)

	config := config{
		conn: 		conn,
		bulbsMap: 	bulbs,
	}
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

