package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (app *App) respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		app.Logger.Error(fmt.Sprintf("%s %s", msg, err.Error()))
	}

	if code > 499 {
		app.Logger.Warn("Responding with 5XX error")
	}

	type errorResponse struct {
		Error string `json:"error"`
	}
	app.respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func (app *App) respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Add("Content-Type", "application/json")

	data, err := json.Marshal(payload)
	if err != nil {
		app.Logger.Error(fmt.Sprintf("Could not marshal the response: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}
