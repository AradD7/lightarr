package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AradD7/lightarr/internal/database"
	"github.com/AradD7/lightarr/internal/wiz"
)

func (app *App) handlerGetBulbs(w http.ResponseWriter, r *http.Request) {
	var bulbs []*wiz.Bulb
	for _, bulb := range app.BulbsMap {
		bulbs = append(bulbs, bulb)
	}
	app.respondWithJSON(w, http.StatusOK, bulbs)
}

func (app *App) handlerUpdateBulbName(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Mac  string `json:"mac"`
		Name string `json:"name"`
	}

	var params parameters
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		app.respondWithError(w, http.StatusBadRequest, "failed to read json data", err)
		return
	}

	if err := app.Db.UpdateBulbName(r.Context(), database.UpdateBulbNameParams{
		Mac:  params.Mac,
		Name: params.Name,
	}); err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Failed to find bulb with given mac", err)
		return
	}
	app.BulbsMap[params.Mac].Name = params.Name

	app.respondWithJSON(w, http.StatusOK, "Updated!")
}

func (app *App) handlerFlashBulb(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Mac string `json:"mac"`
	}

	var params parameters
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		app.respondWithError(w, http.StatusBadRequest, "failed to read json data", err)
		return
	}

	bulb, ok := app.BulbsMap[params.Mac]
	if !ok {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Mac", nil)
		return
	}

	go bulb.Flash(app.Conn)
}

func (app *App) handlerRefreshBulbs(w http.ResponseWriter, r *http.Request) {
	type Resp struct {
		NumNewBulbs int `json:"num"`
	}

	var resp Resp
	app.BulbsMap, resp.NumNewBulbs = wiz.UpdateBulbs(app.Conn, app.BulbsMap, app.Db, app.Logger)
	app.respondWithJSON(w, http.StatusOK, resp)
}

func (app *App) handlerUpdateBulbType(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Mac  string `json:"mac"`
		Type string `json:"type"`
	}

	var params parameters
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		app.respondWithError(w, http.StatusBadRequest, "failed to read json data", err)
		return
	}

	if err := app.Db.UpdateBulbType(r.Context(), database.UpdateBulbTypeParams{
		Mac: params.Mac,
	}); err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Failed to find bulb with given mac", err)
		return
	}
	app.BulbsMap[params.Mac].Type = params.Type

	app.respondWithJSON(w, http.StatusOK, "Updated!")
}
