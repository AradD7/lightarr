package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
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
			Id:        int(acc.ID),
			Title:     acc.Title,
			Thumbnail: acc.Thumb.String,
		})
	}
	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *config) handlerGetAllDevices(w http.ResponseWriter, r *http.Request) {
	players, err := cfg.db.GetAllDevices(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get players from database", err)
		return
	}
	var resp []PlexDevice
	for _, player := range players {
		resp = append(resp, PlexDevice{
			Id:      player.ID,
			Name:    player.Name,
			Product: player.Product,
		})
	}
	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *config) handlerAddAccount(w http.ResponseWriter, r *http.Request) {
	var account PlexAccount
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&account); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode json", err)
		return
	}

	if _, err := cfg.db.AddPlexAccount(r.Context(), database.AddPlexAccountParams{
		ID:    int64(account.Id),
		Title: account.Title,
		Thumb: sql.NullString{
			Valid:  true,
			String: account.Thumbnail,
		},
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to add account to database", err)
		return
	}
	respondWithJSON(w, http.StatusOK, "Account added!")
}

func (cfg *config) handlerAddDevice(w http.ResponseWriter, r *http.Request) {
	var player PlexDevice
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&player); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode json", err)
		return
	}

	if _, err := cfg.db.AddPlexDevice(r.Context(), database.AddPlexDeviceParams{
		ID:      player.Id,
		Name:    player.Name,
		Product: player.Product,
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to add device to database", err)
		return
	}
	respondWithJSON(w, http.StatusOK, "Added device")
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

func (cfg *config) handlerDeleteDevice(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("deviceId")
	if err := cfg.db.DeleteDevice(r.Context(), id); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Found no player with %s id", id), err)
		return
	}
	respondWithJSON(w, http.StatusOK, "Deleted!")
}

func (cfg *config) handlerPlexAllAccounts(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("GET", "https://clients.plex.tv/api/home/users", nil)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create a GET request", err)
		return
	}
	req.Header.Set("X-Plex-Token", cfg.plexToken)
	req.Header.Set("X-Plex-Client-Identifier", cfg.plexClientId)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "GET request failed", err)
		return
	}
	defer resp.Body.Close()

	type User struct {
		Id    string `xml:"id,attr"`
		Title string `xml:"title,attr"`
		Thumb string `xml:"thumb,attr"`
	}
	type result struct {
		XMLName xml.Name `xml:"MediaContainer"`
		Users   []User   `xml:"User"`
	}
	var res result
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&res)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode the XML response", err)
		return
	}

	var accounts []PlexAccount
	for _, user := range res.Users {
		userId, _ := strconv.Atoi(user.Id)
		accounts = append(accounts, PlexAccount{
			Id:        userId,
			Title:     user.Title,
			Thumbnail: user.Thumb,
		})
	}
	respondWithJSON(w, http.StatusOK, accounts)
}

func (cfg *config) handlerPlexAllDevices(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("GET", "https://clients.plex.tv/api/v2/devices", nil)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create a GET request", err)
		return
	}
	req.Header.Set("X-Plex-Token", cfg.plexToken)
	req.Header.Set("X-Plex-Client-Identifier", cfg.plexClientId)
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "GET request failed", err)
		return
	}
	defer resp.Body.Close()

	type parameters struct {
		ClientIdentifier string `json:"clientIdentifier"`
		Name             string `json:"name"`
		Product          string `json:"product"`
	}

	var params []parameters
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode the json response", err)
		return
	}

	var res []PlexDevice
	for _, device := range params {
		res = append(res, PlexDevice{
			Id:      device.ClientIdentifier,
			Name:    device.Name,
			Product: device.Product,
		})
	}

	respondWithJSON(w, http.StatusOK, res)
}
