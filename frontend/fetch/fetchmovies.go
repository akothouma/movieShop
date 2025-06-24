package fetch

import (
	"encoding/json"
	"fmt"
	"os"

	"movieshop/backend/internals/models"
	"movieshop/frontend/fetch"
)

type MoviesResponse struct {
	Page         int             `json:"page"`
	Results      []models.Movies `json:"results"`
	TotalPages   int             `json:"total_pages"`
	TotalResults int             `json:"total_results"`
}

func BasicMovie() ([]models.Movies, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("TMDB_API_KEY not found in environment variables")
	}

	url := "https://api.themoviedb.org/3/discover/movie?include_adult=true&include_video=true&language=en-US&page=1&sort_by=popularity.desc"

	headers := map[string]string{
		"accept":        "application/json",
		"Authorization": "Bearer " + apiKey,
	}

	body, err := fetch.Fetch(url, headers)
	if err != nil {
		return nil, fmt.Errorf("error fetching movies: %w", err)
	}

	var moviesResponse MoviesResponse
	if err := json.Unmarshal(body, &moviesResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return moviesResponse.Results, nil
}
