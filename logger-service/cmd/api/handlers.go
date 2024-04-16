package main

import (
	"fmt"
	"log-service/cmd/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLogs(w http.ResponseWriter, r *http.Request) {
	var requestPayload JSONPayload

	_ = app.ReadJSON(w, r, &requestPayload)

	fmt.Println("requestPayload: FE", requestPayload)
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	fmt.Println("err", err)
	if err != nil {
		_ = app.ErrorJSON(w, err)
		return
	}

	resp := JsonResponse{
		Error:   false,
		Message: "logged",
	}

	app.WriteJSON(w, http.StatusAccepted, resp)

}
