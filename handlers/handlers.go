package handlers

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"main.go/models"
	"main.go/utils"
)

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

func HandleCSVUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	jobID := uuid.New().String()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, "Failed to read CSV record", http.StatusInternalServerError)
			return
		}
		companyRevenue, err := utils.CleanAndConvertToNumeric(record[5])
		if err != nil {
			http.Error(w, "Invalid company revenue format", http.StatusInternalServerError)
			return
		}

		creditAmount, err := utils.CleanAndConvertToNumeric(record[6])
		if err != nil {
			http.Error(w, "Invalid credit amount format", http.StatusInternalServerError)
			return
		}

		query := `INSERT INTO clients (id, name, phone_number, country, gender, company, company_revenue, credit_amount) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
		_, err = db.Exec(query, jobID, record[0], record[1], record[2], record[3], record[4], companyRevenue, creditAmount)
		if err != nil {
			log.Printf("Failed to insert into database: %v", err)
			http.Error(w, "Failed to insert into database", http.StatusInternalServerError)
			return
		}
	}

	jobData := map[string]interface{}{"jobID": jobID, "status": "success"}
	jobDataBytes, _ := json.Marshal(jobData)

	_, err = http.Post("http://localhost:8080/process-job", "application/json", bytes.NewReader(jobDataBytes))
	if err != nil {
		http.Error(w, "Failed to proceed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("CSV processed successfully"))
}

func HandleJobProcessing(w http.ResponseWriter, r *http.Request) {
	log.Println("Processing job request")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		log.Println("Invalid request method")
		return
	}

	var requestBody map[string]string
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v", err)
		return
	}

	jobID, exists := requestBody["jobID"]
	if !exists {
		http.Error(w, "Missing jobID", http.StatusBadRequest)
		log.Println("Missing jobID in request")
		return
	}

	codes, err := models.FetchCallingCodes(db)
	if err != nil {
		http.Error(w, "Failed to fetch calling codes", http.StatusInternalServerError)
		log.Printf("Failed to fetch calling codes: %v", err)
		return
	}

	clients, err := models.FetchClients(db, jobID)
	if err != nil {
		http.Error(w, "Failed to fetch clients", http.StatusInternalServerError)
		return
	}

	validClients := models.ValidateClients(clients, codes)
	if len(validClients) == 0 {
		log.Println("No valid clients found")
		http.Error(w, "No valid clients found", http.StatusInternalServerError)
		return
	}

	for _, client := range validClients {
		if client.CreditAmount > 2000000 {
			log.Printf("Processing client %s with high credit amount", client.ID)
			links, err := utils.SearchGoogle(client.Company, client.Country)
			if err != nil {
				http.Error(w, "Failed to search Google", http.StatusInternalServerError)
				log.Printf("Failed to search Google for client %s: %v", client.ID, err)
				return
			}

			if err := models.SaveClientToDatabase(db, client); err != nil {
				http.Error(w, "Failed to save client to database", http.StatusInternalServerError)
				log.Printf("Failed to save client %s to database: %v", client.ID, err)
				return
			}

			if err := models.SaveLinksToDatabase(db, client, links); err != nil {
				http.Error(w, "Failed to save links to database", http.StatusInternalServerError)
				log.Printf("Failed to save links for client %s: %v", client.ID, err)
				return
			}

		} else {
			log.Printf("Processing client %s with low credit amount - BELOW", client.ID)

			if err := models.SaveClientToDatabase(db, client); err != nil {
				http.Error(w, "Failed to save client to database", http.StatusInternalServerError)
				log.Printf("Failed to save client %s to database: %v", client.ID, err)
				return
			}
		}
	}

	log.Println("Job processed successfully")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Job processed successfully"))
}
