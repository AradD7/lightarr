package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AradD7/lightarr/internal/database"
)

func (cfg *config) handlerGetAllAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := cfg.db.GetAllAccounts(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get accounts from database", err)
		return
	}
	var resp []PlexAccount
	for _, acc := range accounts {
		resp = append(resp, PlexAccount{
			Id: 	int(acc.ID),
			Title:  acc.Title,
		})
	}
	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *config) handlerGetAllPlayers(w http.ResponseWriter, r *http.Request) {
	players, err := cfg.db.GetAllPlayers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get players from database", err)
		return
	}
	var resp []PlexPlayer
	for _, player := range players {
		resp = append(resp, PlexPlayer{
			Uuid: 		player.Uuid,
			Title:  	player.Title,
			PublicAddr: player.PublicAddress,
		})
	}
	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *config) handlerAddAccounts(w http.ResponseWriter, r *http.Request) {
	var accounts []PlexAccount
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&accounts); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode json", err)
		return
	}

	addedAccounts := 0
	for _, account := range accounts{
		if _, err := cfg.db.AddPlexAccount(r.Context(), database.AddPlexAccountParams{
			ID: 	int64(account.Id),
			Title:  account.Title,
		}); err != nil {
			addedAccounts -= 1
		}
		addedAccounts += 1
	}
	respondWithJSON(w, http.StatusOK, fmt.Sprintf("Added %d/%d of the accounts", addedAccounts, len(accounts)))
}

func (cfg *config) handlerAddPlayers(w http.ResponseWriter, r *http.Request) {
	var players []PlexPlayer
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&players); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode json", err)
		return
	}

	addedPlayers := 0
	for _, player := range players {
		if _, err := cfg.db.AddPlexPlayer(r.Context(), database.AddPlexPlayerParams{
			Uuid:  			player.Uuid,
			Title: 			player.Title,
			PublicAddress:  player.PublicAddr,
		}); err != nil {
			addedPlayers -= 1
		}
		addedPlayers += 1
	}
	respondWithJSON(w, http.StatusOK, fmt.Sprintf("Added %d/%d of the players", addedPlayers, len(players)))
}


func (cfg *config) handlerDeleteAccount(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("accountId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Id", err)
		return
	}
	if err := cfg.db.DeleteAccount(r.Context(), int64(id)); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Found no account with %d id", id), err)
		return
	}
	respondWithJSON(w, http.StatusOK, "Deleted!")
}

func (cfg *config) handlerDeletePlayer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("playerId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Id", err)
		return
	}
	if err := cfg.db.DeleteAccount(r.Context(), int64(id)); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Found no player with %d id", id), err)
		return
	}
	respondWithJSON(w, http.StatusOK, "Deleted!")
}

func (cfg *config) handlerPlexWebhook(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to parse form", err)
		return
	}

	payload := r.FormValue("payload")
	if payload == "" {
		respondWithError(w, http.StatusBadRequest, "No payload found", nil)
		return
	}

	var params PlexPayload
	if err := json.Unmarshal([]byte(payload), &params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode JSON", err)
		return
	}

	action, ruleId := cfg.triggersRule(params)
	if action == nil {
		return
	}
	cfg.executeActions(action)
	fmt.Printf("Rule %s was triggered\n", ruleId)
}

func (cfg *config) handlerPlexAllAccounts(w http.ResponseWriter, r *http.Request) {
	//todo
}

func (cfg *config) handlerPlexAllPlayers(w http.ResponseWriter, r *http.Request) {
	//todo
}
