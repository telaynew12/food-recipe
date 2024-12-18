package login

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

//var db *sql.DB

// JWT secret key
var jwtSecret = []byte("Te$%la%ynew2580A!mba!chew5870Msg&ie@")

// User represents a user record in the database
type User struct {
	ID         string
	Email      string
	Password   string
	IsVerified bool
	Role       string
}

// LoginInput represents the structure of the login input
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginOutput represents the structure of the login response
type LoginOutput struct {
	AccessToken string `json:"accessToken"`
	UserID      int    `json:"userId"`
	Role        string `json:"role"`
}

func GenerateJWT(userID string, role string) (string, error) {
	claims := jwt.MapClaims{
		"userId": userID, // Use string userID
		"role":   role,
		"exp":    time.Now().Add(time.Hour * 1).Unix(), // 1-hour expiration
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateCredentials(db *sql.DB, email, password string) (*User, error) {
	var user User
	query := `SELECT id, email, password, is_verified, role FROM users WHERE email = $1 AND is_verified = true`

	// Log the query and parameter
	log.Printf("Executing query: %s with email: %s", query, email)

	err := db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password, &user.IsVerified, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No user found or user not verified for email: %s", email)
			return nil, errors.New("invalid email or user not verified")
		}
		log.Printf("Database error: %v", err)
		return nil, err
	}

	// Log user details (hide sensitive info)
	log.Printf("User retrieved: ID=%s, Email=%s, IsVerified=%t, Role=%s", user.ID, user.Email, user.IsVerified, user.Role)

	// Compare hashed passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Printf("Password mismatch for user ID: %s", user.ID)
		return nil, errors.New("invalid password")
	}

	log.Printf("User authenticated: ID=%s, Email=%s", user.ID, user.Email)
	return &user, nil
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			log.Printf("Failed to read request body: %v", err)
			http.Error(w, `{"message": "Invalid request payload"}`, http.StatusBadRequest)
			return
		}

		// Parse email and password
		var email, password string
		if inputData, ok := requestBody["input"].(map[string]interface{}); ok {
			if nestedInput, ok := inputData["input"].(map[string]interface{}); ok {
				if credentials, ok := nestedInput["credentials"].(map[string]interface{}); ok {
					email, _ = credentials["email"].(string)
					password, _ = credentials["password"].(string)
				}
			}
		}

		// Validate email and password
		if email == "" || password == "" {
			log.Printf("Missing email or password: email=%s, password=%s", email, password)
			http.Error(w, `{"message": "Email and password are required"}`, http.StatusBadRequest)
			return
		}

		// Validate user credentials
		user, err := ValidateCredentials(db, email, password)
		if err != nil {
			log.Printf("Login failed for email %s: %v", email, err)
			http.Error(w, `{"message": "Invalid credentials"}`, http.StatusUnauthorized)
			return
		}

		// Generate JWT token
		accessToken, err := GenerateJWT(user.ID, user.Role)
		if err != nil {
			log.Printf("Failed to generate JWT: %v", err)
			http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
			return
		}

		// Respond with the generated token
		response := map[string]interface{}{
			"accessToken": accessToken,
			"userId":      user.ID,
			"role":        user.Role,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Failed to encode response: %v", err)
			http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
		}
	}
}

// func Login(){
// 		// Connect to database
// 		err := connectToDatabase()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer db.Close()
// }
