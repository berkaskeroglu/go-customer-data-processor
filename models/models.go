package models

import (
	"database/sql"
	"fmt"
	"log"
)

type Client struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	PhoneNumber    string  `json:"phone_number"`
	Country        string  `json:"country"`
	Gender         string  `json:"gender"`
	Company        string  `json:"company"`
	CompanyRevenue float64 `json:"company_revenue"`
	CreditAmount   float64 `json:"credit_amount"`
}

func FetchCallingCodes(db *sql.DB) (map[string]string, error) {
	rows, err := db.Query("SELECT country_name, phone_code FROM calling_codes")
	if err != nil {
		return nil, fmt.Errorf("database query error: %v", err)
	}
	defer rows.Close()

	codes := make(map[string]string)
	for rows.Next() {
		var country, code string
		if err := rows.Scan(&country, &code); err != nil {
			return nil, fmt.Errorf("row scan error: %v", err)
		}
		codes[country] = code
	}

	return codes, nil
}

func FetchClients(db *sql.DB, jobID string) ([]Client, error) {
	rows, err := db.Query("SELECT id, name, phone_number, country, gender, company, company_revenue, credit_amount FROM clients WHERE id = $1", jobID)
	if err != nil {
		return nil, fmt.Errorf("database query error: %v", err)
	}
	defer rows.Close()

	var clients []Client
	for rows.Next() {
		var client Client
		if err := rows.Scan(&client.ID, &client.Name, &client.PhoneNumber, &client.Country, &client.Gender, &client.Company, &client.CompanyRevenue, &client.CreditAmount); err != nil {
			return nil, fmt.Errorf("row scan error: %v", err)
		}
		clients = append(clients, client)
	}

	return clients, nil
}

func ValidateClients(clients []Client, codes map[string]string) []Client {
	var validClients []Client

	for _, client := range clients {
		code, exists := codes[client.Country]
		if !exists {
			log.Printf("No matching phone code for country %s", client.Country)
			continue
		}

		if client.PhoneNumber == code {
			validClients = append(validClients, client)
		} else {
			log.Printf("Phone number %s does not match country code %s for client %s", client.PhoneNumber, code, client.ID)
		}
	}
	log.Printf("Number of valid clients: %d", len(validClients))
	return validClients
}

func SaveClientToDatabase(db *sql.DB, client Client) error {
	query := `INSERT INTO verified_clients (id, name, phone_number, country, gender, company, company_revenue, credit_amount)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	result, err := db.Exec(query, client.ID, client.Name, client.PhoneNumber, client.Country, client.Gender, client.Company, client.CompanyRevenue, client.CreditAmount)
	if err != nil {
		return fmt.Errorf("failed to insert client into database: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not get rows affected: %w", err)
	}

	log.Printf("Inserted client %s, rows affected: %d", client.ID, rowsAffected)
	return nil
}

func SaveLinksToDatabase(db *sql.DB, client Client, links []string) error {
	query := `INSERT INTO links (name, url) VALUES ($1, $2)`

	for _, link := range links {
		_, err := db.Exec(query, client.Name, link)
		if err != nil {
			return fmt.Errorf("failed to insert link for client %s: %w", client.Name, err)
		}
	}
	return nil
}
