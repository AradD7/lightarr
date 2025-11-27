package main

import (
	"net/http"
)

func (cfg *config) handlerGetAllAccounts(w http.ResponseWriter, r *http.Request) {
	_, err := cfg.db.GetAllAccounts(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get accounts from database", err)
		return
	}
}

func (cfg *config) handlerGetAllPlayers(w http.ResponseWriter, r *http.Request) {
	_, err := cfg.db.GetAllPlayers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get players from database", err)
		return
	}
}

func (cfg *config) handlerAddAccounts(w http.ResponseWriter, r *http.Request) {
	//todo
}

func (cfg *config) handlerAddDevices(w http.ResponseWriter, r *http.Request) {
	//todo
}


func (cfg *config) handlerDeleteAccount(w http.ResponseWriter, r *http.Request) {
	//todo
}

func (cfg *config) handlerDeleteDevice(w http.ResponseWriter, r *http.Request) {
	//todo
}

