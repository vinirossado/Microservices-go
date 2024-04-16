package main

import (
	"fmt"
	"log"
	"net/http"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	fmt.Println("MailService: requestPayload: ", requestPayload)

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		log.Println("MailService: Error after SendSMTPMessage", err)
		app.errorJSON(w, err)
		return
	}

	payload := JsonResponse{
		Error:   false,
		Message: "sent to " + requestPayload.To,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
