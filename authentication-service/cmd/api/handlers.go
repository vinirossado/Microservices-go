package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.ReadJSON(w, r, &requestPayload)
	if err != nil {
		err := app.ErrorJSON(w, err, http.StatusBadRequest)
		if err != nil {
			return
		}
		return
	}

	fmt.Println("Request Payload", requestPayload)
	user, err := app.Models.User.GetByEmail(requestPayload.Email)

	if err != nil {
		err := app.ErrorJSON(w, errors.New("invalid Credentials"), http.StatusBadRequest)
		if err != nil {
			return
		}
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)

	if err != nil || !valid {
		err := app.ErrorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		if err != nil {
			return
		}
		return
	}

	err = app.logRequest("authentication", fmt.Sprintf("User %s logged in", user.Email))

	if err != nil {
		_ = app.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	err = app.WriteJSON(w, http.StatusAccepted, payload)
	if err != nil {
		return
	}

}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	req, err := http.NewRequest(http.MethodPost, logServiceURL, bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
