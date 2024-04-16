package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
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
		err := app.ErrorJSON(w, err, http.StatusBadRequest)
		if err != nil {
			return
		}
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.logItem(w, requestPayload.Log)
	default:
		err := app.ErrorJSON(w, err, http.StatusBadRequest)
		if err != nil {
			return
		}
	}

}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, err := json.MarshalIndent(entry, "", "\t")

	fmt.Println("jsonData", string(jsonData))

	if err != nil {
		_ = app.ErrorJSON(w, err)
		return

	}

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))

	if err != nil {
		err := app.ErrorJSON(w, err)
		if err != nil {
			return
		}

		request.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		response, err := client.Do(request)

		fmt.Println("RESPONSE", response)
		fmt.Println("RESPONSE ERR FROM LOGGER", err)
		if err != nil {
			return
		}

		defer response.Body.Close()

		if response.StatusCode != http.StatusAccepted {
			err := app.ErrorJSON(w, errors.New("error calling log service"), http.StatusInternalServerError)
			if err != nil {
				return
			}
			return
		}

		var payload JsonResponse
		payload.Error = false
		payload.Message = "Logged"

		err = app.WriteJSON(w, http.StatusAccepted, payload)
		if err != nil {
			return
		}

	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))

	if err != nil {
		err := app.ErrorJSON(w, err)
		if err != nil {
			return
		}
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		err := app.ErrorJSON(w, err)
		if err != nil {
			return
		}
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	if response.StatusCode == http.StatusUnauthorized {
		err := app.ErrorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		if err != nil {
			return
		}
		return
	} else if response.StatusCode != http.StatusAccepted {
		err := app.ErrorJSON(w, errors.New("error calling auth service"), http.StatusInternalServerError)
		if err != nil {
			return
		}
		return
	}

	var authResponse JsonResponse

	err = json.NewDecoder(response.Body).Decode(&authResponse)

	if err != nil {
		err := app.ErrorJSON(w, err)
		if err != nil {
			return
		}
		return
	}

	if authResponse.Error {
		err := app.ErrorJSON(w, errors.New(authResponse.Message), http.StatusUnauthorized)
		if err != nil {
			return
		}
		return
	}

	var payload JsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = authResponse.Data

	_ = app.WriteJSON(w, http.StatusAccepted, payload)
}
