package main

import (
	"github.com/AradD7/lightarr/internal/wiz"
)

type PlexAccount struct {
	Id		int 	`json:"id"`
	Title	string 	`json:"title"`
}

type PlexPlayer struct {
	PublicAddr 	string `json:"publicAddress"`
	Title 		string `json:"title"`
	Uuid 		string `json:"uuid"`
}

type PlexPayload struct {
	Event 	string 		`json:"event"`
	Account PlexAccount `json:"Account"`
	Player 	PlexPlayer	`json:"Player"`
}

type WizAction struct {
	Command 	wiz.WizCommand	`json:"command"`
	BulbId		[]string		`json:"BulbIds"`
}

type Rule struct {
	Condition 	PlexPayload `json:"condition"`
	Action 		[]WizAction	`json:"action"`
}
