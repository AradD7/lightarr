package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *config) handlerPrintRule(w http.ResponseWriter, r *http.Request) {
	var rule Rule
	decoder := json.NewDecoder(r.Body)
	log.Println(r.Body)
	if err := decoder.Decode(&rule); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to unmarshal json", err)
		return
	}
	log.Println(rule)
	respondWithJSON(w, http.StatusOK, "Rule Printed!")
}

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

	for idx, device := range rule.Condition.Device {
		dev, err := cfg.db.GetPlexDeviceById(r.Context(), device.Id)
		if err != nil {
			log.Printf("%s is not a valid device ID. Check saved devices or simply restart the app", dev.ID)
			continue
		}
		rule.Condition.Device[idx].Name = dev.Name
		rule.Condition.Device[idx].Product = dev.Product
	}

	for idx, account := range rule.Condition.Account {
		acc, err := cfg.db.GetPlexAccountById(r.Context(), int64(account.Id))
		if err != nil {
			log.Printf("%d is not a valid account ID. Check saved accounts or simply restart the app", acc.ID)
			continue
		}
		rule.Condition.Account[idx].Title = acc.Title
		rule.Condition.Account[idx].Thumbnail = acc.Thumb.String
	}

	err := cfg.addRule(rule.Condition.Event, rule.Condition.Account, rule.Condition.Device, rule.Action)
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
