# Customer Data Processor

This project processes customer data from CSV files, validates it based on country phone codes, and searches for company information on Google for customers with high credit limits. The results are saved to a database.

## Features

- Upload CSV files and save the data to the database.
- Validate customer data based on country phone codes.
- Search for company information on Google for customers with high credit limits.
- Save search results to the database.

## Installation

Follow the steps below to set up and run the project on your local machine.

### Prerequisites

- [Go 1.20 or higher](https://golang.org/dl/)
- [PostgreSQL](https://www.postgresql.org/download/)
- Google Custom Search API key and CX value (obtainable from [Google Cloud Console](https://console.cloud.google.com/)).

### Steps

1. Clone the repository:
   ```bash
   git clone https://github.com/berkaskeroglu/go-customer-data-processor.git
   cd go-customer-data-processor
2. Install the dependencies:
   ```bash
   go mod download
3. Run the project:
   ```bash
   go run main.go

### Usage

1. To upload a CSV file, use the following command:
   ```bash
   curl -F "file=@path/to/your/file.csv" http://localhost:8080/upload-csv
2. To start processing, use the following command:
   ```bash
   curl -X POST http://localhost:8080/process-job -d '{"jobID": "your-job-id"}'
    
