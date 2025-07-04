package handler

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"

	"gopkg.in/gomail.v2"
)

// RequestBody defines the structure of the incoming JSON request.
type RequestBody struct {
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	Message   string `json:"message"` // This will be the HTML content of the email body
}

// EmailTemplateData holds the data for the HTML email template.
type EmailTemplateData struct {
	Message template.HTML // Use template.HTML to prevent escaping of HTML tags
}

// Handler is the main entry point for the Vercel serverless function.
func Handler(w http.ResponseWriter, r *http.Request) {
	// --- CORS and Basic Setup ---
	allowedOrigins := map[string]bool{
		"https://tbilisi.hackclub.com":    true,
		"https://tbilisihc.andrinoff.com": true,
	}

	// 2. Get the origin of the request.
	origin := r.Header.Get("Origin")

	// 3. Check if the origin is in your whitelist and set the header.
	if allowedOrigins[origin] {
		w.Header().Set("Access-control-Allow-Origin", origin)
	}

	// The rest of the CORS and handler logic remains the same.
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight CORS requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// (The rest of your existing code continues from here...)

	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// --- Load Credentials ---
	senderEmail := os.Getenv("YAHOO_EMAIL")
	senderPassword := os.Getenv("YAHOO_APP_PASSWORD")

	if senderEmail == "" || senderPassword == "" {
		log.Println("Error: Environment variables for email not set.")
		http.Error(w, `{"error": "Server configuration error."}`, http.StatusInternalServerError)
		return
	}

	// --- Parse Request ---
	var reqBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, `{"error": "Invalid request body."}`, http.StatusBadRequest)
		return
	}

	if reqBody.Recipient == "" || reqBody.Subject == "" || reqBody.Message == "" {
		http.Error(w, `{"error": "recipient, subject, and message are required fields."}`, http.StatusBadRequest)
		return
	}

	// --- Compose and Send HTML Email ---
	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail)
	m.SetHeader("To", reqBody.Recipient)
	m.SetHeader("Subject", reqBody.Subject)

	// Generate HTML body from the generic template
	htmlBody, err := generateHTML(reqBody.Message)
	if err != nil {
		log.Printf("Failed to generate HTML email body: %v", err)
		http.Error(w, `{"error": "Failed to generate email content."}`, http.StatusInternalServerError)
		return
	}

	m.SetBody("text/html", htmlBody)

	// --- Send Email ---
	d := gomail.NewDialer("smtp.mail.yahoo.com", 587, senderEmail, senderPassword)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email to %s: %v", reqBody.Recipient, err)
		http.Error(w, `{"error": "Failed to send email."}`, http.StatusInternalServerError)
		return
	}

	// --- Success Response ---
	log.Printf("Email sent successfully to %s", reqBody.Recipient)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Email sent successfully!"})
}

// generateHTML creates the email's HTML content from a generic template.
func generateHTML(messageContent string) (string, error) {
	const tpl = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; background-color: #f4f4f7; margin: 0; padding: 0; color: #333; }
        .container { max-width: 600px; margin: 20px auto; background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 4px 15px rgba(0,0,0,0.1); }
        .header { background: linear-gradient(
            45deg,
            rgb(10, 74, 89) 0%,
            rgb(11, 82, 124) 8.333%,
            rgb(13, 92, 159) 16.667%,
            rgb(16, 104, 194) 25%,
            rgb(18, 118, 227) 33.333%,
            rgb(21, 133, 255) 41.667%,
            rgb(24, 150, 255) 50%,
            rgb(28, 169, 255) 58.333%,
            rgb(31, 188, 255) 66.667%,
            rgb(35, 209, 255) 75%,
            rgb(39, 230, 255) 83.333%,
            rgb(43, 251, 255) 91.667%,
            rgb(47, 255, 255) 100%
        ); padding: 40px; text-align: center; }
        .header img { max-width: 150px; }
        .content { padding: 40px; line-height: 1.6; }
        .footer { text-align: center; padding: 20px; font-size: 12px; color: #888; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
       	<svg
            width="150"
            height="60"
            viewBox="0 0 150 60"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
        >
            <style>
                .logo-text {
                    font-family: -apple-system, BlinkMacSystemFont,
                        "Segoe UI", Roboto, Helvetica, Arial,
                        sans-serif, "Apple Color Emoji",
                        "Segoe UI Emoji";
                    font-size: 45px;
                    font-weight: 500;
                    fill: white;
                    text-anchor: middle;
                    dominant-baseline: central;
                }
            </style>
            <text x="50%" y="50%" class="logo-text">thc</text>
        </svg>
        </div>
        <div class="content">
            {{.Message}}
        </div>
        <div class="footer">
            <p>&copy; 2025 Tbilisi Hack Club</p>
        </div>
    </div>
</body>
</html>`

	t, err := template.New("email").Parse(tpl)
	if err != nil {
		return "", err
	}

	// Pass the message content to the template
	data := EmailTemplateData{Message: template.HTML(messageContent)}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
