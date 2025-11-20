package wiz

import (
	"encoding/json"
	"net"
)

type Bulb struct {
	Id 	 string	`json:"id"`
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
