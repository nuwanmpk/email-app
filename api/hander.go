package api

import (
	"email-app/internal/email"
	"encoding/json"
	"net/http"
)

type ContactForm struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func SendEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var form ContactForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	emailBody := "Name: " + form.Name + "\n" +
		"Email: " + form.Email + "\n\n" +
		form.Message

	if err := email.SendEmail(form.Subject, emailBody); err != nil {
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email sent successfully"))
}
