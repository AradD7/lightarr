package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/AradD7/lightarr/internal/database"
	"github.com/AradD7/lightarr/internal/wiz"
)

type GetPilotParams struct {
	Result struct {
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
			Mac:  bulb.Mac,
			Ip:   net.ParseIP(bulb.Ip),
			Name: bulb.Name,
			Addr: &net.UDPAddr{
				IP:   net.ParseIP(bulb.Ip),
				Port: 38899,
			},
			IsReachable: false,
			Type:        bulb.Type,
		}

		bulbsMap[bulb.Mac] = &currentBulb
	}

	cfg.UpdateBulbs(conn, bulbsMap)
}

func (cfg *config) UpdateBulbs(conn *net.UDPConn, bulbsMap map[string]*wiz.Bulb) int {
	if bulbsMap != nil {
		fmt.Println("Updating current bulbs and checking for additional light bulbs on the network...")
	} else {
		fmt.Println("No cache file found. Discovering light bulbs on the network...")
	}

	subnet := os.Getenv("WIZ_SUBNET")
	if subnet == "" {
		subnet = "192.168.1.0/24"
	}

	discoveredBulbs := cfg.scanSubnetForBulbs(conn, subnet)

	bulbId := 0
	updatedBulbs := 0

	for mac, discoveredBulb := range discoveredBulbs {
		if cachedBulb, ok := bulbsMap[mac]; ok {
			bulbsMap[mac].IsReachable = true
			if !cachedBulb.Ip.Equal(discoveredBulb.Ip) {
				bulbsMap[mac].Addr.IP = discoveredBulb.Ip
				bulbsMap[mac].Ip = discoveredBulb.Ip
				err := cfg.db.UpdateBulbIp(context.Background(), database.UpdateBulbIpParams{
					Mac:       mac,
					Ip:        discoveredBulb.Ip.String(),
					UpdatedAt: time.Now(),
				})
				if err != nil {
					fmt.Printf("Failed to update DB: %v\n", err.Error())
				}
				updatedBulbs += 1
			}
		} else {
			bulbsMap[mac] = discoveredBulb
			_, err := cfg.db.AddBulb(context.Background(), database.AddBulbParams{
				Mac:         mac,
				Ip:          discoveredBulb.Ip.String(),
				Name:        "WizBulb",
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

func (cfg *config) scanSubnetForBulbs(conn *net.UDPConn, subnet string) map[string]*wiz.Bulb {
	bulbs := make(map[string]*wiz.Bulb)
	var mu sync.Mutex
	var wg sync.WaitGroup

	ip, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		fmt.Printf("Invalid subnet %s: %v\n", subnet, err)
		return bulbs
	}

	fmt.Printf("Scanning subnet %s for Wiz bulbs...\n", subnet)

	semaphore := make(chan struct{}, 50)

	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); incrementIP(ip) {
		// Skip network and broadcast addresses
		if ip[3] == 0 || ip[3] == 255 {
			continue
		}

		wg.Add(1)
		semaphore <- struct{}{}

		go func(targetIP string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			if bulb := cfg.probeBulb(conn, targetIP); bulb != nil {
				mu.Lock()
				bulbs[bulb.Mac] = bulb
				mu.Unlock()
			}
		}(ip.String())
	}

	wg.Wait()
	fmt.Printf("Scan complete. Found %d bulbs.\n", len(bulbs))
	return bulbs
}

func (cfg *config) probeBulb(conn *net.UDPConn, ip string) *wiz.Bulb {
	targetAddr := &net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: 38899,
	}

	getPilotMsg := []byte(`{"method":"getPilot","params":{}}`)
	_, err := conn.WriteToUDP(getPilotMsg, targetAddr)
	if err != nil {
		return nil
	}

	conn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))

	buffer := make([]byte, 1024)
	n, remoteAddr, err := conn.ReadFromUDP(buffer)
	if err != nil || n == 0 {
		return nil
	}

	var params GetPilotParams
	if err := json.Unmarshal(buffer[:n], &params); err != nil {
		return nil
	}

	if params.Result.Mac == "" {
		return nil
	}

	return &wiz.Bulb{
		Ip:   remoteAddr.IP,
		Name: "WizBulb",
		Mac:  params.Result.Mac,
		Addr: &net.UDPAddr{
			IP:   remoteAddr.IP,
			Port: 38899,
		},
		IsReachable: true,
		Type:        "normal",
	}
}

func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
