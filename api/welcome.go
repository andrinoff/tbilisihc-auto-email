package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"gopkg.in/gomail.v2"
)

// RequestBody defines the structure of the incoming JSON request.
type RequestBody struct {
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	Message   string `json:"message"`
}

// Handler is the main entry point for the Vercel serverless function.
func Handler(w http.ResponseWriter, r *http.Request) {
	allowedOrigin := "https://tbilisi.hackclub.com"
	w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight CORS requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	// --- 1. Basic Setup & Security ---

	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// --- 2. Load Credentials from Environment Variables ---
	senderEmail := os.Getenv("YAHOO_EMAIL")
	senderPassword := os.Getenv("YAHOO_APP_PASSWORD")

	if senderEmail == "" || senderPassword == "" {
		log.Println("Error: YAHOO_EMAIL or YAHOO_APP_PASSWORD environment variables not set.")
		http.Error(w, `{"error": "Server configuration error."}`, http.StatusInternalServerError)
		return
	}

	// --- 3. Parse Incoming Request ---
	var reqBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, `{"error": "Invalid request body."}`, http.StatusBadRequest)
		return
	}

	// Basic validation
	if reqBody.Recipient == "" || reqBody.Subject == "" || reqBody.Message == "" {
		http.Error(w, `{"error": "recipient, subject, and message are required."}`, http.StatusBadRequest)
		return
	}

	// --- 4. Compose and Send Email using gomail ---
	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail)
	m.SetHeader("To", reqBody.Recipient)
	m.SetHeader("Subject", reqBody.Subject)
	m.SetBody("text/plain", reqBody.Message)

	// Yahoo's SMTP server details
	d := gomail.NewDialer("smtp.mail.yahoo.com", 587, senderEmail, senderPassword)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email to %s: %v", reqBody.Recipient, err)
		http.Error(w, `{"error": "Failed to send email."}`, http.StatusInternalServerError)
		return
	}

	// --- 5. Send Success Response ---
	log.Printf("Email sent successfully to %s", reqBody.Recipient)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Email sent successfully!"})
}
