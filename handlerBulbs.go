package main

import (
	"net/http"

	"github.com/AradD7/lightarr/internal/wiz"
)

func (cfg *config) handlerGetBulbs(w http.ResponseWriter, r *http.Request) {
	var bulbs []*wiz.Bulb
	for _, bulb := range cfg.bulbsMap {
		bulbs = append(bulbs, bulb)
	}
	respondWithJSON(w, http.StatusOK, bulbs)
}
