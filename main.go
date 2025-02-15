package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"main.go/handlers"
)

const (
	dbUser     = "postgres"
	dbPassword = "1234"
	dbName     = "postgres"
)

var db *sql.DB

func main() {
	var err error
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName)
	db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/upload-csv", handlers.HandleCSVUpload)
	http.HandleFunc("/process-job", handlers.HandleJobProcessing)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
