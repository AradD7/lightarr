package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AradD7/lightarr/internal/config"
	"github.com/AradD7/lightarr/internal/database"
)

func (cfg *config.Config) handlerGetAllRules(w http.ResponseWriter, r *http.Request) {
	cfg.respondWithJSON(w, http.StatusOK, cfg.rules)
}

func (cfg *config.Config) handlerAddRule(w http.ResponseWriter, r *http.Request) {
	var rule Rule
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rule); err != nil {
		cfg.respondWithError(w, http.StatusBadRequest, "Failed to unmarshal json", err)
		return
	}

	for idx, device := range rule.Condition.Device {
		dev, err := cfg.db.GetPlexDeviceById(r.Context(), device.Id)
		if err != nil {
			cfg.logger.Warn(fmt.Sprint("%s is not a valid device ID. Check saved devices or simply restart the app", dev.ID))
			continue
		}
		rule.Condition.Device[idx].Name = dev.Name
		rule.Condition.Device[idx].Product = dev.Product
	}

	for idx, account := range rule.Condition.Account {
		acc, err := cfg.db.GetPlexAccountById(r.Context(), int64(account.Id))
		if err != nil {
			cfg.logger.Warn(fmt.Sprint("%d is not a valid account ID. Check saved accounts or simply restart the app", acc.ID))
			continue
		}
		rule.Condition.Account[idx].Title = acc.Title
		rule.Condition.Account[idx].Thumbnail = acc.Thumb.String
	}

	err := cfg.addRule(rule.Condition.Event, rule.Condition.Account, rule.Condition.Device, rule.Action)
	if err != nil {
		cfg.respondWithError(w, http.StatusInternalServerError, "Failed to add rule", err)
		return
	}

	cfg.respondWithJSON(w, http.StatusOK, "Rule Added!")
}

func (cfg *config.Config) handlerDeleteRule(w http.ResponseWriter, r *http.Request) {
	ruleId := r.PathValue("ruleId")
	err := cfg.deleteRule(ruleId)
	if err != nil {
		cfg.respondWithError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	cfg.respondWithJSON(w, http.StatusOK, "Deleted!")
}

func (cfg *config.Config) handlerUpdateRuleName(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}

	var params parameters
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		cfg.respondWithError(w, http.StatusBadRequest, "failed to read json data", err)
		return
	}

	err := cfg.db.UpdateRuleName(r.Context(), database.UpdateRuleNameParams{
		ID: params.Id,
		Name: sql.NullString{
			Valid:  true,
			String: params.Name,
		},
	})

	if err != nil {
		cfg.respondWithError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	for idx, rule := range cfg.rules {
		if rule.Id == params.Id {
			cfg.rules[idx].Name = params.Name
			cfg.respondWithJSON(w, http.StatusOK, "Updated!")
			return
		}
	}

	cfg.respondWithError(w, http.StatusInternalServerError, "Something went wrong with changing the rule name. Please restart the server", nil)
}
