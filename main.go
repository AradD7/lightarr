package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"github.com/AradD7/lightarr/internal/wiz"
)

func main() {
	conn, _ := net.ListenUDP("udp", &net.UDPAddr{Port: 0})
	fmt.Printf("Opened a UDP connection on %s\n", conn.LocalAddr().String())
	defer conn.Close()

	bulbs := loadBulbs(conn)

	cmd := wiz.NewSetPilotWizCommand()
	cmd.TurnOn()

	for _, bulb := range bulbs {
		bulb.Execute(conn, cmd)
	}
}

func loadBulbs(conn *net.UDPConn) []wiz.Bulb {
	var bulbs []wiz.Bulb
	cacheMap := make(map[string]*wiz.Bulb)
	nameMap := make(map[string]string)

	data, err := os.ReadFile(wiz.BulbCacheFile)
	if err != nil {
		fmt.Printf("Could not open cache file: %s\n", err.Error())
		return wiz.DiscoverBulbs(conn, nil, nil)
	}

	if err = json.Unmarshal(data, &bulbs); err != nil {
		fmt.Printf("Failed to unmarshal the data in cache file: %s\n", err.Error())
		return wiz.DiscoverBulbs(conn, nil, nil)
	}

	for _, bulb := range bulbs {
		cacheMap[bulb.Ip.String()] = &bulb
		nameMap[bulb.Mac] = bulb.Name
	}
	return wiz.DiscoverBulbs(conn, cacheMap, nameMap)
}
