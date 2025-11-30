package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/AradD7/lightarr/internal/database"
	"github.com/AradD7/lightarr/internal/wiz"
)

type GetPilotParams struct {
	Result struct{
		Mac string `json:"mac"`
	} `json:"result"`
}

func (cfg *config) LoadBulbs(conn *net.UDPConn) {
	bulbsMap := make(map[string]*wiz.Bulb)

	data, err := cfg.db.GetAllBulbs(context.Background())
	if err != nil {
		fmt.Printf("Could not read bulbs from db: %s\n", err.Error())
		return
	}

	for _, bulb := range data {
		currentBulb := wiz.Bulb{
			Mac: 		 bulb.Mac,
			Ip: 		 net.ParseIP(bulb.Ip),
			Name: 		 bulb.Name,
			Addr: 		 &net.UDPAddr{
				IP: 	net.ParseIP(bulb.Ip),
				Port: 	38899,
			},
			IsReachable: false,
		}

		bulbsMap[bulb.Mac] = &currentBulb
	}

	cfg.UpdateBulbs(conn, bulbsMap)
}

func (cfg *config) UpdateBulbs(conn *net.UDPConn, bulbsMap map[string]*wiz.Bulb) int {
	buffer := make([]byte, 1024)

	if bulbsMap != nil {
		fmt.Println("Updating current bulbs and checking for additional light bulbs on the network...")
	} else {
		fmt.Println("No cahce file found. Discovering light bulbs on the network...")
	}

	broadcastAddr := &net.UDPAddr{
		IP: net.IPv4bcast,
		Port: 38899,
	}
	conn.WriteToUDP([]byte(`{"method":"getPilot"}`), broadcastAddr)

	bulbId := 0
	updatedBulbs := 0
	var params GetPilotParams
	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if n == 0 || err != nil {
			break
		}

		json.Unmarshal(buffer[:n], &params)
		if cachedBulb, ok := bulbsMap[params.Result.Mac]; ok {
			bulbsMap[params.Result.Mac].IsReachable = true
			if !cachedBulb.Ip.Equal(remoteAddr.IP) {
				bulbsMap[params.Result.Mac].Addr.IP = remoteAddr.IP
				bulbsMap[params.Result.Mac].Ip = remoteAddr.IP
				err := cfg.db.UpdateBulbIp(context.Background(), database.UpdateBulbIpParams{
					Mac: 		params.Result.Mac,
					Ip: 		remoteAddr.IP.String(),
					UpdatedAt: 	time.Now(),
				})
				if err != nil {
					fmt.Printf("Failed to update DB: %v\n", err.Error())
				}
				updatedBulbs += 1
			}
		} else {
			bulbsMap[params.Result.Mac] = &wiz.Bulb{
				Ip: 		 remoteAddr.IP,
				Name: 		 "WizBulb",
				Mac: 		 params.Result.Mac,
				Addr: 		 &net.UDPAddr{
					IP: 	remoteAddr.IP,
					Port: 	38899,
				},
				IsReachable: true,
			}
			_, err := cfg.db.AddBulb(context.Background(), database.AddBulbParams{
				Mac: 		 params.Result.Mac,
				Ip: 		 remoteAddr.IP.String(),
				Name: 		 "WizBulb",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				IsReachable: true,
			})
			if err != nil {
				fmt.Printf("Failed to add bulb to DB: %v\n", err.Error())
			}
			bulbId += 1
		}
	}

	switch updatedBulbs {
	case 0:
		break
	case 1:
		fmt.Printf("Updated %d bulb\n", updatedBulbs)
	default:
		fmt.Printf("Updated %d bulbs\n", updatedBulbs)
	}
	switch bulbId {
	case 0:
		fmt.Println("Found no new bulbs.")
	case 1:
		fmt.Printf("Found %d new bulb\n", bulbId)
	default:
		fmt.Printf("Found %d new bulbs\n", bulbId)
	}

	cfg.bulbsMap = bulbsMap
	return bulbId
}
