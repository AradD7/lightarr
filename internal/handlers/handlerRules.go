package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AradD7/lightarr/internal/database"
	"github.com/AradD7/lightarr/internal/rules"
)

func (app *App) HandlerGetAllRules(w http.ResponseWriter, r *http.Request) {
	app.respondWithJSON(w, http.StatusOK, app.Rules)
}

func (app *App) HandlerAddRule(w http.ResponseWriter, r *http.Request) {
	var rule rules.Rule
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rule); err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Failed to unmarshal json", err)
		return
	}

	for idx, device := range rule.Condition.Device {
		dev, err := app.Db.GetPlexDeviceById(r.Context(), device.Id)
		if err != nil {
			app.Logger.Warn(fmt.Sprint("%s is not a valid device ID. Check saved devices or simply restart the app", dev.ID))
			continue
		}
		rule.Condition.Device[idx].Name = dev.Name
		rule.Condition.Device[idx].Product = dev.Product
	}

	for idx, account := range rule.Condition.Account {
		acc, err := app.Db.GetPlexAccountById(r.Context(), int64(account.Id))
		if err != nil {
			app.Logger.Warn(fmt.Sprint("%d is not a valid account ID. Check saved accounts or simply restart the app", acc.ID))
			continue
		}
		rule.Condition.Account[idx].Title = acc.Title
		rule.Condition.Account[idx].Thumbnail = acc.Thumb.String
	}

	var err error
	app.Rules, err = rules.AddRule(app.Rules, rule.Condition.Event, rule.Condition.Account, rule.Condition.Device, rule.Action, app.Db)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to add rule", err)
		return
	}

	app.respondWithJSON(w, http.StatusOK, "Rule Added!")
}

func (app *App) HandlerDeleteRule(w http.ResponseWriter, r *http.Request) {
	ruleId := r.PathValue("ruleId")
	var err error
	app.Rules, err = rules.DeleteRule(app.Rules, ruleId, app.Db)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	app.respondWithJSON(w, http.StatusOK, "Deleted!")
}

func (app *App) HandlerUpdateRuleName(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}

	var params parameters
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		app.respondWithError(w, http.StatusBadRequest, "failed to read json data", err)
		return
	}

	err := app.Db.UpdateRuleName(r.Context(), database.UpdateRuleNameParams{
		ID: params.Id,
		Name: sql.NullString{
			Valid:  true,
			String: params.Name,
		},
	})

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	for idx, rule := range app.Rules {
		if rule.Id == params.Id {
			app.Rules[idx].Name = params.Name
			app.respondWithJSON(w, http.StatusOK, "Updated!")
			return
		}
	}

	app.respondWithError(w, http.StatusInternalServerError, "Something went wrong with changing the rule name. Please restart the server", nil)
}
