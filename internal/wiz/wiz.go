package wiz

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

type Bulb struct {
	Ip 	 net.IP `json:"ip"`
	Name string `json:"name"`
	Mac  string `json:"mac"`
	Addr *net.UDPAddr `json:"-"`
}

type Color struct {
	R, G, B int
}

type WizCommand struct {
	Method string 	 `json:"method"`
	Params WizParams `json:"params"`
}

type WizParams struct {
	State 	*bool `json:"state,omitempty"`
	Dimming *int  `json:"dimming,omitempty"`
	R 		*int  `json:"r,omitempty"`
	G 		*int  `json:"g,omitempty"`
	B 		*int  `json:"b,omitempty"`
	Temp 	*int  `json:"temp,omitempty"`
	SceneId *int  `json:"sceneId,omitempty"`
}

type GetPilotParams struct {
	Result struct{
		Mac string `json:"mac"`
	} `json:"result"`
}

const BulbCacheFile = "bulbs.json"

func LoadBulbs(conn *net.UDPConn) map[string]*Bulb {
	var bulbs []Bulb
	bulbsMap := make(map[string]*Bulb)

	data, err := os.ReadFile(BulbCacheFile)
	if err != nil {
		fmt.Printf("Could not open cache file: %s\n", err.Error())
		return UpdateBulbs(conn, bulbsMap)
	}

	if err = json.Unmarshal(data, &bulbs); err != nil {
		fmt.Printf("Failed to unmarshal the data in cache file: %s\n", err.Error())
		return UpdateBulbs(conn, bulbsMap)
	}

	for _, bulb := range bulbs {
		bulb.Addr = &net.UDPAddr{
			IP: 	bulb.Ip,
			Port: 	38899,
		}
		bulbsMap[bulb.Mac] = &bulb
	}
	return UpdateBulbs(conn, bulbsMap)
}

func UpdateBulbs(conn *net.UDPConn, bulbsMap map[string]*Bulb) map[string]*Bulb {
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
	var params GetPilotParams
	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if n == 0 || err != nil {
			break
		}

		json.Unmarshal(buffer[:n], &params)
		if cachedBulb, ok := bulbsMap[params.Result.Mac]; ok {
			if !cachedBulb.Ip.Equal(remoteAddr.IP) {
				bulbsMap[params.Result.Mac].Addr.IP = remoteAddr.IP
				bulbsMap[params.Result.Mac].Ip = remoteAddr.IP
			}
		} else {
			bulbsMap[params.Result.Mac] = &Bulb{
				Ip: 	remoteAddr.IP,
				Name: 	"WizBulb",
				Mac: 	params.Result.Mac,
				Addr: 	&net.UDPAddr{
					IP: 	remoteAddr.IP,
					Port: 	38899,
				},
			}
			bulbId += 1
		}
	}

	switch bulbId {
	case 0:
		fmt.Println("Found no new bulbs.")
	case 1:
		fmt.Printf("Found %d new bulb\n", bulbId)
	default:
		fmt.Printf("Found %d new bulbs\n", bulbId)
	}

	var bulbs []Bulb
	for _, bulb := range bulbsMap {
		bulbs = append(bulbs, *bulb)
	}
	data, err := json.MarshalIndent(bulbs, "", "   ")
	if err != nil {
		fmt.Printf("Failed to marshal bulbs to cache: %s\n", err.Error())
		return bulbsMap
}
	err = os.WriteFile(BulbCacheFile, data, 0644)
	if err != nil {
		fmt.Printf("Failed to cache the bulbs: %s\n", err.Error())
	}

	return bulbsMap
}

func (b *Bulb) Execute(conn *net.UDPConn, cmd WizCommand) error {
	data, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	_, err = conn.WriteToUDP(data, b.Addr)
	return err
}

func NewSetPilotWizCommand() WizCommand {
	return WizCommand{
		Method: "setPilot",
		Params: WizParams{},
	}
}

func NewGetPilotWizCommand() WizCommand {
	return WizCommand{
		Method: "getPilot",
		Params: WizParams{},
	}
}

func (w *WizCommand) SetDim(brightness int) {
	w.Params.Dimming = &brightness
}

func (w *WizCommand) SetColor(c Color) {
	w.Params.R = &c.R
	w.Params.G = &c.G
	w.Params.B = &c.B
}

func (w *WizCommand) TurnOff() {
	off := false
	w.Params.State = &off
}

func (w *WizCommand) TurnOn() {
	on := true
	w.Params.State = &on
}

func (w *WizCommand) SetTemp(temp int) {
	w.Params.Temp = &temp
}
