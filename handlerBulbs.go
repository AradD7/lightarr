package main

import (
	"encoding/json"
	"net/http"

	"github.com/AradD7/lightarr/internal/database"
	"github.com/AradD7/lightarr/internal/wiz"
)

func (cfg *config) handlerGetBulbs(w http.ResponseWriter, r *http.Request) {
	var bulbs []*wiz.Bulb
	for _, bulb := range cfg.bulbsMap {
		bulbs = append(bulbs, bulb)
	}
	respondWithJSON(w, http.StatusOK, bulbs)
}

func (cfg *config) handlerUpdateBulbName(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Mac  string `json:"mac"`
		Name string `json:"name"`
	}

	var params parameters
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to read json data", err)
		return
	}

	if err := cfg.db.UpdateBulbName(r.Context(), database.UpdateBulbNameParams{
		Mac:  params.Mac,
		Name: params.Name,
	}); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to find bulb with given mac", err)
		return
	}
	cfg.bulbsMap[params.Mac].Name = params.Name

	respondWithJSON(w, http.StatusOK, "Updated!")
}

func (cfg *config) handlerFlashBulb(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Mac string `json:"mac"`
	}

	var params parameters
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to read json data", err)
		return
	}

	bulb, ok := cfg.bulbsMap[params.Mac]
	if !ok {
		respondWithError(w, http.StatusBadRequest, "Invalid Mac", nil)
		return
	}

	go bulb.Flash(cfg.conn)
}

func (cfg *config) handlerRefreshBulbs(w http.ResponseWriter, r *http.Request) {
	type Resp struct {
		NumNewBulbs int `json:"num"`
	}

	respondWithJSON(w, http.StatusOK, Resp{
		NumNewBulbs: cfg.UpdateBulbs(cfg.conn, cfg.bulbsMap),
	})
}

func (cfg *config) handlerUpdateBulbType(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Mac  string `json:"mac"`
		Type string `json:"type"`
	}

	var params parameters
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to read json data", err)
		return
	}

	if err := cfg.db.UpdateBulbType(r.Context(), database.UpdateBulbTypeParams{
		Mac: params.Mac,
	}); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to find bulb with given mac", err)
		return
	}
	cfg.bulbsMap[params.Mac].Type = params.Type

	respondWithJSON(w, http.StatusOK, "Updated!")
}
