package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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
