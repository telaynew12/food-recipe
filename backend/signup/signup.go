

package signup

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"time"

	//"golang.org/x/crypto/bcrypt" // Import bcrypt library

	_ "github.com/lib/pq" // PostgreSQL driver
	//"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Input struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	} `json:"input"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}

type VerifyRequest struct {
	Input struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	} `json:"input"`
}
type VerifyResponse struct {
	Message string `json:"message"`
}

func generateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}

func sendEmail(email, verificationCode string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	smtpUser := "telaynew11@gmail.com"
	smtpPass := "kerpgqurcvdgzdtu"

	from := smtpUser
	to := []string{email}
	subject := "Subject: Verification Code\n"
	body := fmt.Sprintf("Your verification code is: %s", verificationCode)
	msg := []byte(subject + "\n" + body)

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
	if err != nil {
		log.Printf("Failed to send email: %v\n", err)
		return err
	}

	log.Printf("Verification email sent to %s successfully.\n", email)
	return nil
}

func hashPassword(password string) (string, error) {
	// Generate a bcrypt hash of the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func RegisterHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var requestBody map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err == nil {
		requestBodyBytes, _ := json.Marshal(requestBody)
		log.Printf("Incoming request body: %s", string(requestBodyBytes))
	} else {
		log.Printf("Failed to read request body for logging: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var email, name, password string
	if inputData, ok := requestBody["input"].(map[string]interface{}); ok {
		if input, ok := inputData["input"].(map[string]interface{}); ok {
			email, _ = input["email"].(string)
			name, _ = input["name"].(string)
			password, _ = input["password"].(string)
		}
	}

	if email == "" || name == "" || password == "" {
		http.Error(w, "Email, name, and password are required", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := hashPassword(password)
	if err != nil {
		log.Printf("Password hashing error: %v\n", err)
		http.Error(w, `{"message":"Failed to process password"}`, http.StatusInternalServerError)
		return
	}

	verificationCode := generateVerificationCode()

	var isVerified bool
	var createdAt time.Time
	query := `SELECT is_verified, created_at FROM users WHERE email = $1`
	err = db.QueryRow(query, email).Scan(&isVerified, &createdAt)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Database error: %v\n", err)
		http.Error(w, `{"message":"Database error"}`, http.StatusInternalServerError)
		return
	}
	if isVerified {
		response := RegisterResponse{Message: "User already verified. No further action needed."}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if err == nil && !isVerified {
		updateQuery := `UPDATE users SET verification_code = $1, created_at = $2 WHERE email = $3`
		_, err = db.Exec(updateQuery, verificationCode, time.Now(), email)
		if err != nil {
			log.Printf("Database update error: %v\n", err)
			http.Error(w, `{"message":"Failed to update verification code"}`, http.StatusInternalServerError)
			return
		}
	} else if err == sql.ErrNoRows {
		insertQuery := `INSERT INTO users (email, name, password, verification_code, created_at, is_verified) VALUES ($1, $2, $3, $4, $5, false)`
		_, err = db.Exec(insertQuery, email, name, hashedPassword, verificationCode, time.Now())
		if err != nil {
			log.Printf("Database insert error: %v\n", err)
			http.Error(w, `{"message":"Failed to save user in database"}`, http.StatusInternalServerError)
			return
		}
	}

	if err := sendEmail(email, verificationCode); err != nil {
		log.Printf("Email sending error: %v\n", err)
		http.Error(w, `{"message":"Failed to send verification email"}`, http.StatusInternalServerError)
		return
	}

	response := RegisterResponse{Message: "Verification code sent successfully."}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func VerifyHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	// Decode the incoming request body
	var requestBody map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	log.Printf("Incoming request body: %v", requestBody)

	// Extract the email and code
	var email, code string
	if inputData, ok := requestBody["input"].(map[string]interface{}); ok {
		if input, ok := inputData["input"].(map[string]interface{}); ok {
			email, _ = input["email"].(string)
			code, _ = input["code"].(string)
		}
	}

	// Validate email and code
	if email == "" || code == "" {
		log.Println("Email or code not provided in the request.")
		http.Error(w, "Email or code not provided", http.StatusBadRequest)
		return
	}

	var dbCode string
	var createdAt time.Time
	query := `SELECT verification_code, created_at FROM users WHERE email = $1 AND is_verified = false`
	err := db.QueryRow(query, email).Scan(&dbCode, &createdAt)

	if err == sql.ErrNoRows || dbCode != code {
		http.Error(w, `{"message":"Invalid verification code or email"}`, http.StatusBadRequest)
		return
	} else if err != nil {
		log.Printf("Database error: %v\n", err)
		http.Error(w, `{"message":"Database error"}`, http.StatusInternalServerError)
		return
	}

	// Check if the code is older than 5 minutes.
	if time.Since(createdAt) > 5*time.Minute {
		http.Error(w, `{"message":"Verification code expired"}`, http.StatusBadRequest)
		return
	}

	// Update the user to set is_verified to true.
	updateQuery := `UPDATE users SET is_verified = true WHERE email = $1`
	_, err = db.Exec(updateQuery, email)
	if err != nil {
		log.Printf("Database update error: %v\n", err)
		http.Error(w, `{"message":"Failed to update verification status"}`, http.StatusInternalServerError)
		return
	}

	response := VerifyResponse{Message: "Email verified successfully."}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
