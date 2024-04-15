package main

import (
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
		err := app.ErrorJSON(w, errors.New("Invalid Credentials"), http.StatusBadRequest)
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
