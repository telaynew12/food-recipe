package main

import (
    "backend/chapa"
    "backend/image"
    "backend/login"
    "backend/notify"
    "backend/signup"
    "database/sql"
    "log"
    "net/http"
    "fmt"
)

func connectToDatabase() (*sql.DB, error) {
    dbHost := "localhost"
    dbPort := "5431"
    dbUser := "telay"
    dbPassword := "123456"
    dbName := "uog"

    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("database connection error: %w", err)
    }
    return db, nil
}

func main() {
    db, err := connectToDatabase()
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

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

    // Set up the routes for Chapa payment
    http.HandleFunc("/payment/callback", func(w http.ResponseWriter, r *http.Request) {
        chapa.PaymentCallback(w, r, db)
    })
    http.HandleFunc("/start-payment", chapa.StartPaymentHandler)

    log.Println("Starting server on http://localhost:8085")
    err = http.ListenAndServe(":8085", nil)
    if err != nil {
        log.Fatal(err)
    }
}
