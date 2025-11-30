package wiz

import (
	"encoding/json"
	"net"
	"time"
)

type Bulb struct {
	Ip 	 		net.IP 		 `json:"ip"`
	Name 		string 		 `json:"name"`
	Mac  		string 		 `json:"mac"`
	Addr 		*net.UDPAddr `json:"-"`
	IsReachable bool 		 `json:"isReachable"`
	Type 		string 		 `json:"type"`
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

func (b *Bulb) drainBuffer(conn *net.UDPConn) {
    conn.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
    buf := make([]byte, 1024)
    for {
        _, _, err := conn.ReadFromUDP(buf)
        if err != nil {
            break
        }
    }
}

func (b *Bulb) Flash(conn *net.UDPConn) {
	type parameters struct {
		Result struct {
			Dimming int `json:"dimming"`
		} `json:"result"`
	}
	var params parameters
	buffer := make([]byte, 1024)
	b.drainBuffer(conn)
	getPilot := NewGetPilotWizCommand()
	b.Execute(conn, getPilot)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, _, _ := conn.ReadFromUDP(buffer)
	json.Unmarshal(buffer[:n], &params)

	flashCmd := NewSetPilotWizCommand()
	flashCmd.SetDim(0)
	b.Execute(conn, flashCmd)
	time.Sleep(800 * time.Millisecond)
	flashCmd.SetDim(100)
	b.Execute(conn, flashCmd)
	time.Sleep(800 * time.Millisecond)
	flashCmd.SetDim(0)
	b.Execute(conn, flashCmd)
	time.Sleep(800 * time.Millisecond)
	flashCmd.SetDim(params.Result.Dimming)
	b.Execute(conn, flashCmd)

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
