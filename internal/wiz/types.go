package wiz

import (
	"net"
)

type Bulb struct {
	Ip          net.IP       `json:"ip"`
	Name        string       `json:"name"`
	Mac         string       `json:"mac"`
	Addr        *net.UDPAddr `json:"-"`
	IsReachable bool         `json:"isReachable"`
	Type        string       `json:"type"`
}

type Color struct {
	R, G, B int
}

type WizCommand struct {
	Method string    `json:"method"`
	Params WizParams `json:"params"`
}

type WizParams struct {
	State   *bool `json:"state,omitempty"`
	Dimming *int  `json:"dimming,omitempty"`
	R       *int  `json:"r,omitempty"`
	G       *int  `json:"g,omitempty"`
	B       *int  `json:"b,omitempty"`
	Temp    *int  `json:"temp,omitempty"`
	SceneId *int  `json:"sceneId,omitempty"`
}

type WizAction struct {
	Command  WizCommand `json:"command"`
	BulbsMac []string   `json:"bulbsMac"`
}
