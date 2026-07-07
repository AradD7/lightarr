package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (cfg *config) handlerPlexWebhook(w http.ResponseWriter, r *http.Request) {
	cfg.logger.Debug("Recieved Plex Payload")
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		cfg.respondWithError(w, http.StatusBadRequest, "Failed to parse form", err)
		return
	}

	payload := r.FormValue("payload")
	if payload == "" {
		cfg.respondWithError(w, http.StatusBadRequest, "No payload found", nil)
		return
	}

	var params PlexPayload
	if err := json.Unmarshal([]byte(payload), &params); err != nil {
		cfg.respondWithError(w, http.StatusBadRequest, "Failed to decode JSON", err)
		return
	}

	action, ruleId := cfg.triggersRule(params)
	if action == nil {
		return
	}
	cfg.executeActions(action)
	cfg.logger.Info(fmt.Sprintf("Rule %s was triggered", ruleId))
}
