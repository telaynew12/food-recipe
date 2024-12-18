package notify

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/gomail.v2"
)

// EventPayload represents the incoming payload structure
type EventPayload struct {
	Event struct {
		Data struct {
			New struct {
				Email      string `json:"email"`
				IsVerified bool   `json:"is_verified"`
			} `json:"new"`
			Old struct {
				IsVerified bool `json:"is_verified"`
			} `json:"old"`
		} `json:"data"`
	} `json:"event"`
}

// sendEmailnotify sends an email notification
func sendEmailnotify(email string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "telaynew11@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Registration Successful")
	m.SetBody("text/plain", "Congratulations! Your registration is now verified.")

	// Replace "kerp gqur cvdg zdtu" with your app password or appropriate credentials
	d := gomail.NewDialer("smtp.gmail.com", 587, "telaynew11@gmail.com", "kerp gqur cvdg zdtu")

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

// HandleEvent processes incoming events and sends email if verified
func HandleEvent(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var payload EventPayload

	// Parse JSON payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Println("Error parsing payload:", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	newData := payload.Event.Data.New
	oldData := payload.Event.Data.Old

	// Check if the is_verified status changed to true (and wasn't already true before)
	if oldData.IsVerified != newData.IsVerified && newData.IsVerified {
		if err := sendEmailnotify(newData.Email); err != nil {
			log.Println("Error sending email:", err)
			http.Error(w, "Failed to send email", http.StatusInternalServerError)
			return
		}
		log.Println("Email sent successfully to:", newData.Email)
	}

	// Optionally, you can use the db here if needed (for example, to log event details in the database)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "processed"})
}
