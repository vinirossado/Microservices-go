package main

import (
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {

	payload := JsonResponse{
		Error:   false,
		Message: "Hit the Broker",
	}

	_ = app.WriteJSON(w, http.StatusOK, payload)
}
