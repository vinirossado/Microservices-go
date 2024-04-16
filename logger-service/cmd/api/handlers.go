package main

import (
	"log-service/cmd/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLogs(w http.ResponseWriter, r *http.Request) {
	var requestPayload JSONPayload

	_ = app.readJSON(w, r, &requestPayload)

	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Models.LogEntry.Insert(event)

	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	resp := JsonResponse{
		Error:   false,
		Message: "logged",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}
