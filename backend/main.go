package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq" // PostgreSQL driver
)

//var db *sql.DB

func main() {
	// Database connection details
	dbHost := "localhost"
	dbPort := "5439"
	dbUser := "telay"
	dbPassword := "123456"
	dbName := "uog"

	// Connection string (DSN)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	// Initialize database connection
	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v\n", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Database connection error: %v\n", err)
	}

	// Define HTTP routes
	http.HandleFunc("/events", handleEvent)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/verify", verifyHandler)
	http.HandleFunc("/login", LoginHandler(db))

	// // Paths to your SSL certificate and key files
	// certFile := "C:\\Users\\telay\\Desktop\\my-fullstack-app\\backend\\server.crt"
	// keyFile := "C:\\Users\\telay\\Desktop\\my-fullstack-app\\backend\\server.key"

	// Start the HTTPS server
	// err = http.ListenAndServeTLS(":8085")

	log.Println("Starting server on http://localhost:8085")
	err = http.ListenAndServe(":8085", nil) // nil uses the default multiplexer (http.ServeMux)
	if err != nil {
		log.Fatal(err)
	}
}
