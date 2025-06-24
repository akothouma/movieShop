package fetch

import (
	"encoding/json"
	"fmt"
	"os"

	"movieshop/backend/internals/models"
	"movieshop/frontend/fetch"
)

func RatingInfo() ([]models.Movies, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	// Get API key from environment
	apiKey := os.Getenv("OMDB_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OMDB_API_KEY not found in environment variables")
	}

	movies := []models.Movies{}

	// Construct URL with API key from environment
	url := fmt.Sprintf("http://www.omdbapi.com/?i=tt3896198&apikey=%s", apiKey)

	body, fetchErr := fetch.Fetch(url)
	if fetchErr != nil {
		return nil, fetchErr
	}

	unmarshal_err := json.Unmarshal(body, &movies)
	if unmarshal_err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", unmarshal_err)
	}
   fmt.Println(movies)
	return movies, nil

}
