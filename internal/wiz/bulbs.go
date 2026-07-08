package wiz

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"os"
	"sync"
	"time"

	"github.com/AradD7/lightarr/internal/database"
)

type GetPilotParams struct {
	Result struct {
		Mac string `json:"mac"`
	} `json:"result"`
}

func LoadBulbs(conn *net.UDPConn, db *database.Queries, logger *slog.Logger) {
	bulbsMap := make(map[string]*Bulb)

	data, err := db.GetAllBulbs(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("Could not read bulbs from db: %s", err.Error()))
		return
	}

	for _, bulb := range data {
		currentBulb := Bulb{
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

	UpdateBulbs(conn, bulbsMap, db, logger)
}

func UpdateBulbs(conn *net.UDPConn, bulbsMap map[string]*Bulb, db *database.Queries, logger *slog.Logger) (map[string]*Bulb, int) {
	if bulbsMap != nil {
		logger.Info("Updating current bulbs and checking for additional light bulbs on the network...")
	} else {
		logger.Info("No cache file found. Discovering light bulbs on the network...")
	}

	subnet := os.Getenv("WIZ_SUBNET")
	if subnet == "" {
		subnet = "192.168.1.0/24"
	}

	discoveredBulbs := scanSubnetForBulbs(conn, subnet, logger)

	bulbId := 0
	updatedBulbs := 0

	for mac, discoveredBulb := range discoveredBulbs {
		if cachedBulb, ok := bulbsMap[mac]; ok {
			bulbsMap[mac].IsReachable = true
			if !cachedBulb.Ip.Equal(discoveredBulb.Ip) {
				bulbsMap[mac].Addr.IP = discoveredBulb.Ip
				bulbsMap[mac].Ip = discoveredBulb.Ip
				err := db.UpdateBulbIp(context.Background(), database.UpdateBulbIpParams{
					Mac:       mac,
					Ip:        discoveredBulb.Ip.String(),
					UpdatedAt: time.Now(),
				})
				if err != nil {
					logger.Error(fmt.Sprintf("Failed to update DB: %v", err.Error()))
				}
				updatedBulbs += 1
			}
		} else {
			bulbsMap[mac] = discoveredBulb
			_, err := db.AddBulb(context.Background(), database.AddBulbParams{
				Mac:         mac,
				Ip:          discoveredBulb.Ip.String(),
				Name:        "WizBulb",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				IsReachable: true,
			})
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to add bulb to DB: %v", err.Error()))
			}
			bulbId += 1
		}
	}

	switch updatedBulbs {
	case 0:
		break
	case 1:
		logger.Info(fmt.Sprintf("Updated %d bulb", updatedBulbs))
	default:
		logger.Info(fmt.Sprintf("Updated %d bulbs", updatedBulbs))
	}

	switch bulbId {
	case 0:
		logger.Info("Found no new bulbs.")
	case 1:
		logger.Info(fmt.Sprintf("Found %d new bulb", bulbId))
	default:
		logger.Info(fmt.Sprintf("Found %d new bulbs", bulbId))
	}

	return bulbsMap, bulbId
}

func scanSubnetForBulbs(conn *net.UDPConn, subnet string, logger *slog.Logger) map[string]*Bulb {
	bulbs := make(map[string]*Bulb)
	var mu sync.Mutex
	var wg sync.WaitGroup

	ip, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		logger.Error(fmt.Sprintf("Invalid subnet %s: %v", subnet, err))
		return bulbs
	}

	logger.Info(fmt.Sprintf("Scanning subnet %s for Wiz bulbs...", subnet))

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

			if bulb := probeBulb(conn, targetIP); bulb != nil {
				mu.Lock()
				bulbs[bulb.Mac] = bulb
				mu.Unlock()
			}
		}(ip.String())
	}

	wg.Wait()
	logger.Info(fmt.Sprintf("Scan complete. Found %d bulbs", len(bulbs)))
	return bulbs
}

func probeBulb(conn *net.UDPConn, ip string) *Bulb {
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

	return &Bulb{
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
