package utils

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"google.golang.org/api/customsearch/v1"
	"google.golang.org/api/option"
)

const (
	gcsAPIKey = "*******"
	gcsCx     = "*******"
)

func CleanAndConvertToNumeric(value string) (float64, error) {
	re := regexp.MustCompile(`[^\d.]`)
	cleaned := re.ReplaceAllString(value, "")

	numericValue, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0, err
	}

	return numericValue, nil
}

func SearchGoogle(company, country string) ([]string, error) {
	query := company + " " + country
	ctx := context.Background()
	svc, err := customsearch.NewService(ctx, option.WithAPIKey(gcsAPIKey))
	if err != nil {
		log.Fatalf("Error creating customsearch service: %v", err)
	}

	resp, err := svc.Cse.List().Cx(gcsCx).Q(query).Do()
	if err != nil {
		log.Fatalf("Error making search request: %v", err)
	}

	var links []string
	for i, item := range resp.Items {
		fmt.Printf("#%d: %s\n", i+1, item.Title)
		fmt.Printf("\t%s\n", item.Snippet)
		links = append(links, item.Link)
	}

	return links, nil
}
