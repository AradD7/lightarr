package main

import (
	"fmt"
	"net"
	"time"
)

type Bulb struct {
	Name string
	Addr *net.UDPAddr
}

func main() {
	var bulbs []Bulb

	conn, _ := net.ListenUDP("udp", &net.UDPAddr{Port: 0})
	defer conn.Close()

	broadcastAddr := &net.UDPAddr{
		IP: net.IPv4bcast,
		Port: 38899,
	}

	conn.WriteToUDP([]byte(`{"method":"getPilot"}`), broadcastAddr)

	buffer := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	bulbId := 1
	for n, remoteAddr, err := conn.ReadFromUDP(buffer); n != 0 && err == nil; n, remoteAddr, err = conn.ReadFromUDP(buffer) {
		fmt.Println("found a bulb:", remoteAddr.IP, remoteAddr.Port)
		bulbs = append(bulbs, Bulb{
			Name: fmt.Sprintf("Bulb%d", bulbId),
			Addr: remoteAddr,
		})
		bulbId += 1
	}

	fmt.Println(bulbs)
}
