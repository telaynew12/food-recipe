package main

import (
	"backend/image"
	"backend/login"
	"backend/notify"
	"backend/signup"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// connectToDatabase establishes a connection to the PostgreSQL database
func connectToDatabase() (*sql.DB, error) {
	dbHost := "localhost"
	dbPort := "5431"
	dbUser := "telay"
	dbPassword := "123456"
	dbName := "uog"

	// Connection string (DSN)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	// Initialize database connection
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test database connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}
	return db, nil
}

func main() {
	// Connect to the database
	db, err := connectToDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// Serve static files from the uploads directory
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	// Initialize routes with database connection
	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		notify.HandleEvent(db, w, r)
	})
	http.HandleFunc("/login", login.LoginHandler(db))
	http.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
		image.FileUploadHandler(db, w, r)
	})
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		signup.RegisterHandler(w, r, db)
	})
	http.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		signup.VerifyHandler(w, r, db)
	})

	// Start the server
	log.Println("Starting server on http://localhost:8085")
	err = http.ListenAndServe(":8085", nil)
	if err != nil {
		log.Fatal(err)
	}
}
