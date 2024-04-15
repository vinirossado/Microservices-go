package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {

	payload := JsonResponse{
		Error:   false,
		Message: "Hit the Broker",
	}

	_ = app.WriteJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {

	var requestPayload RequestPayload

	err := app.ReadJSON(w, r, &requestPayload)

	if err != nil {
		app.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)

	default:
		app.ErrorJSON(w, err, http.StatusBadRequest)
	}

}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {

	jsonData, _ := json.MarshalIndent(a, "", "\t")

	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))

	if err != nil {
		app.ErrorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		app.ErrorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.ErrorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.ErrorJSON(w, errors.New("error calling auth service"), http.StatusInternalServerError)
		return
	}

	var authResponse JsonResponse

	err = json.NewDecoder(response.Body).Decode(&authResponse)

	if err != nil {
		app.ErrorJSON(w, err)
		return
	}

	if authResponse.Error {
		app.ErrorJSON(w, errors.New(authResponse.Message), http.StatusUnauthorized)
		return
	}

	var payload JsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = authResponse.Data

	_ = app.WriteJSON(w, http.StatusAccepted, payload)
}
