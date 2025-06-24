package fetch

import (
	"encoding/json"
	"fmt"
	"movieshop/backend/internals/models"
	"os"

	"github.com/joho/godotenv"
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

	url := "https://api.themoviedb.org/3/discover/movie?include_adult=false&include_video=false&language=en-US&page=1&sort_by=popularity.desc"

	headers := map[string]string{
		"accept":        "application/json",
		"Authorization": "Bearer " + apiKey,
	}

	body, err := Fetch(url, headers)
	if err != nil {
		return nil, fmt.Errorf("error fetching movies: %w", err)
	}

	var moviesResponse MoviesResponse
	if err := json.Unmarshal(body, &moviesResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return moviesResponse.Results, nil
}

type OMDBResponse struct {
	Title    string `json:"Title"`
	Year     string `json:"Year"`
	Rated    string `json:"Rated"`
	Released string `json:"Released"`
	Runtime  string `json:"Runtime"`
	Genre    string `json:"Genre"`
	Director string `json:"Director"`
	Writer   string `json:"Writer"`
	Actors   string `json:"Actors"`
	Plot     string `json:"Plot"`
	Language string `json:"Language"`
	Country  string `json:"Country"`
	Awards   string `json:"Awards"`
	Poster   string `json:"Poster"`
	Ratings  []struct {
		Source string `json:"Source"`
		Value  string `json:"Value"`
	} `json:"Ratings"`
	Metascore  string `json:"Metascore"`
	ImdbRating string `json:"imdbRating"`
	ImdbVotes  string `json:"imdbVotes"`
	ImdbID     string `json:"imdbID"`
	Type       string `json:"Type"`
	Response   string `json:"Response"`
}

func RatingInfo(imdbID string) (*OMDBResponse, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	apiKey := os.Getenv("OMDB_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OMDB_API_KEY not found in environment variables")
	}

	if imdbID == "" {
		return nil, fmt.Errorf("IMDB ID is required")
	}

	url := fmt.Sprintf("http://www.omdbapi.com/?i=%s&apikey=%s", imdbID, apiKey)

	body, err := Fetch(url, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching rating info: %w", err)
	}

	var response OMDBResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if response.Response != "True" {
		return nil, fmt.Errorf("OMDB API returned error response")
	}

	return &response, nil
}

func GetMovieByID(movieID int) (*models.Movies, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("TMDB_API_KEY not found in environment variables")
	}

	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d?language=en-US", movieID)

	headers := map[string]string{
		"accept":        "application/json",
		"Authorization": "Bearer " + apiKey,
	}

	body, err := Fetch(url, headers)
	if err != nil {
		return nil, fmt.Errorf("error fetching movie: %w", err)
	}

	var movie models.Movies
	if err := json.Unmarshal(body, &movie); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &movie, nil
}
