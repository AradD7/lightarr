package handlers

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AradD7/lightarr/internal/database"
	"github.com/AradD7/lightarr/internal/plex"
)

func (app *App) handlerGetAllAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := app.Db.GetAllAccounts(r.Context())
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to get accounts from database", err)
		return
	}
	var resp []plex.PlexAccount
	for _, acc := range accounts {
		resp = append(resp, plex.PlexAccount{
			Id:        int(acc.ID),
			Title:     acc.Title,
			Thumbnail: acc.Thumb.String,
		})
	}
	app.respondWithJSON(w, http.StatusOK, resp)
}

func (app *App) handlerGetAllDevices(w http.ResponseWriter, r *http.Request) {
	players, err := app.Db.GetAllDevices(r.Context())
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to get players from database", err)
		return
	}
	var resp []plex.PlexDevice
	for _, player := range players {
		resp = append(resp, plex.PlexDevice{
			Id:      player.ID,
			Name:    player.Name,
			Product: player.Product,
		})
	}
	app.respondWithJSON(w, http.StatusOK, resp)
}

func (app *App) handlerAddAccount(w http.ResponseWriter, r *http.Request) {
	var account plex.PlexAccount
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&account); err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Failed to decode json", err)
		return
	}

	if _, err := app.Db.AddPlexAccount(r.Context(), database.AddPlexAccountParams{
		ID:    int64(account.Id),
		Title: account.Title,
		Thumb: sql.NullString{
			Valid:  true,
			String: account.Thumbnail,
		},
	}); err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to add account to database", err)
		return
	}
	app.respondWithJSON(w, http.StatusOK, "Account added!")
}

func (app *App) handlerAddDevice(w http.ResponseWriter, r *http.Request) {
	var player plex.PlexDevice
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&player); err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Failed to decode json", err)
		return
	}

	if _, err := app.Db.AddPlexDevice(r.Context(), database.AddPlexDeviceParams{
		ID:      player.Id,
		Name:    player.Name,
		Product: player.Product,
	}); err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to add device to database", err)
		return
	}
	app.respondWithJSON(w, http.StatusOK, "Added device")
}

func (app *App) handlerDeleteAccount(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("accountId"))
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Id", err)
		return
	}
	if err := app.Db.DeleteAccount(r.Context(), int64(id)); err != nil {
		app.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Found no account with %d id", id), err)
		return
	}
	app.respondWithJSON(w, http.StatusOK, "Deleted!")
}

func (app *App) handlerDeleteDevice(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("deviceId")
	if err := app.Db.DeleteDevice(r.Context(), id); err != nil {
		app.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Found no player with %s id", id), err)
		return
	}
	app.respondWithJSON(w, http.StatusOK, "Deleted!")
}

func (app *App) handlerPlexAllAccounts(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("GET", "https://clients.plex.tv/api/home/users", nil)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to create a GET request", err)
		return
	}
	req.Header.Set("X-Plex-Token", app.Cfg.PlexToken)
	req.Header.Set("X-Plex-Client-Identifier", app.Cfg.PlexClientId)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "GET request failed", err)
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
		app.respondWithError(w, http.StatusInternalServerError, "Failed to decode the XML response", err)
		return
	}

	var accounts []plex.PlexAccount
	for _, user := range res.Users {
		userId, _ := strconv.Atoi(user.Id)
		accounts = append(accounts, plex.PlexAccount{
			Id:        userId,
			Title:     user.Title,
			Thumbnail: user.Thumb,
		})
	}
	app.respondWithJSON(w, http.StatusOK, accounts)
}

func (app *App) handlerPlexAllDevices(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("GET", "https://clients.plex.tv/api/v2/devices", nil)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to create a GET request", err)
		return
	}
	req.Header.Set("X-Plex-Token", app.Cfg.PlexToken)
	req.Header.Set("X-Plex-Client-Identifier", app.Cfg.PlexClientId)
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "GET request failed", err)
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
		app.respondWithError(w, http.StatusInternalServerError, "Failed to decode the json response", err)
		return
	}

	var res []plex.PlexDevice
	for _, device := range params {
		res = append(res, plex.PlexDevice{
			Id:      device.ClientIdentifier,
			Name:    device.Name,
			Product: device.Product,
		})
	}

	app.respondWithJSON(w, http.StatusOK, res)
}
