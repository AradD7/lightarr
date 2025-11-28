package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *config) handlerGetAllRules(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, cfg.rules)
}

func (cfg *config) handlerAddRule(w http.ResponseWriter, r *http.Request) {
	var rule Rule
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rule); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to unmarshal json", err)
		return
	}

	err := cfg.addRule(rule.Condition.Event, rule.Condition.Account, rule.Condition.Player, rule.Action)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to add rule", err)
		return
	}

	respondWithJSON(w, http.StatusOK, "Rule Added!")
}

func (cfg *config) handlerDeleteRule(w http.ResponseWriter, r *http.Request) {
	ruleId := r.PathValue("ruleId")
	err := cfg.deleteRule(ruleId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	respondWithJSON(w, http.StatusOK, "Deleted!")
}
