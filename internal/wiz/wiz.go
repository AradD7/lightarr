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

func DiscoverBulbs(conn *net.UDPConn, cacheMap map[string]*Bulb, nameMap map[string]string) []Bulb {
	var bulbs []Bulb
	getPilot := NewGetPilotWizCommand()
	buffer := make([]byte, 1024)

	if cacheMap != nil {
		fmt.Println("Loading cache...")
		for _, val := range cacheMap {
			conn.SetReadDeadline(time.Now().Add(2 * time.Second))
			bulb := Bulb{
				Ip: 	val.Ip,
				Name: 	val.Name,
				Mac: 	val.Mac,
				Addr: 	&net.UDPAddr{
					IP: 	val.Ip,
					Port: 	38899,
				},
			}
			err := bulb.Execute(conn, getPilot)
			if err != nil {
				fmt.Printf("Skipping bulb %s because: %s\n", val.Name, err.Error())
				continue
			}
			n, respAddr, err := conn.ReadFromUDP(buffer)
			if n == 0 || err != nil {
				fmt.Printf("Skipping bulb %s because: %s\n", val.Name, err.Error())
				continue
			}
			if respAddr.IP.Equal(bulb.Ip) {
				bulbs = append(bulbs, bulb)
			}
		}
		fmt.Printf("Loaded %d bulbs from the Cache.\n", len(bulbs))
		fmt.Println("Discovering additional light bulbs on the network...")
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

		if _, ok := cacheMap[remoteAddr.IP.String()]; ok {
			continue
		}
		json.Unmarshal(buffer[:n], &params)
		if name, ok := nameMap[params.Result.Mac]; ok {
			bulbs = append(bulbs, Bulb{
				Ip:   remoteAddr.IP,
				Name: name,
				Mac:  params.Result.Mac,
				Addr: remoteAddr,
			})
		} else {
			bulbs = append(bulbs, Bulb{
				Ip:   remoteAddr.IP,
				Name: "WizBulb",
				Mac:  params.Result.Mac,
				Addr: remoteAddr,
			})
		}
		bulbId += 1
	}

	switch bulbId {
	case 0:
		fmt.Println("Found no new bulbs.")
	case 1:
		fmt.Printf("Found %d new bulb\n", bulbId)
	default:
		fmt.Printf("Found %d new bulbs\n", bulbId)
	}

	data, err := json.MarshalIndent(bulbs, "", "   ")
	if err != nil {
		fmt.Printf("Failed to marshal bulbs to cache: %s\n", err.Error())
		return bulbs
	}
	err = os.WriteFile(BulbCacheFile, data, 0644)
	if err != nil {
		fmt.Printf("Failed to cache the bulbs: %s\n", err.Error())
	}

	return bulbs
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
