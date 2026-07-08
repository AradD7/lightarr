package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AradD7/lightarr/internal/plex"
	"github.com/AradD7/lightarr/internal/rules"
	"github.com/AradD7/lightarr/internal/wiz"
)

func (app *App) HandlerPlexWebhook(w http.ResponseWriter, r *http.Request) {
	app.Logger.Debug("Recieved Plex Payload")
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Failed to parse form", err)
		return
	}

	payload := r.FormValue("payload")
	if payload == "" {
		app.respondWithError(w, http.StatusBadRequest, "No payload found", nil)
		return
	}

	var params plex.PlexPayload
	if err := json.Unmarshal([]byte(payload), &params); err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Failed to decode JSON", err)
		return
	}

	action, ruleId := rules.TriggersRule(app.Rules, params)
	if action == nil {
		return
	}
	wiz.ExecuteActions(app.BulbsMap, action, app.Conn, app.Logger)
	app.Logger.Info(fmt.Sprintf("Rule %s was triggered", ruleId))
}
