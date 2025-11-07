package main

import (
	"fmt"
	"net"
	"github.com/AradD7/lightarr/internal/wiz"
)

func main() {
	conn, _ := net.ListenUDP("udp", &net.UDPAddr{Port: 0})
	fmt.Printf("Opened a UDP connection on %s\n", conn.LocalAddr().String())
	defer conn.Close()

	bulbs := wiz.LoadBulbs(conn)

	config := config{
		conn: 		conn,
		bulbsMap: 	bulbs,
	}

	startRepl(&config)
}

