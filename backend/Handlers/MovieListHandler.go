package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"movieshop/backend/internals/models"
	fetch "movieshop/backend/utils"
)

type CombinedMovieResponse struct {
	models.Movies
	OMDBRating string `json:"omdb_rating,omitempty"`
}

// MovieListHandler returns a handler for the movie list endpoint
func MovieListHandler() http.Handler {
	mux := http.NewServeMux()
	
	// Register the routes
	mux.HandleFunc("/", GetMoviesAPI)
	mux.HandleFunc("/movie", GetMovieDetailsAPI)
	
	return mux
}

// GetMoviesAPI returns movies as JSON
func GetMoviesAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	movies, err := fetch.BasicMovie()
	if err != nil {
		http.Error(w, "Failed to fetch movies", http.StatusInternalServerError)
		return
	}

	if len(movies) == 0 {
		http.Error(w, "No movies found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Add CORS header
	json.NewEncoder(w).Encode(movies)
}

// GetMovieDetailsAPI returns detailed movie info as JSON
func GetMovieDetailsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	movieID, err := strconv.Atoi(idStr)
	if err != nil || movieID < 1 {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	movie, err := fetch.GetMovieByID(movieID)
	if err != nil {
		http.Error(w, "Failed to fetch movie details", http.StatusInternalServerError)
		return
	}

	// Create combined response with OMDB rating if available
	combinedResponse := CombinedMovieResponse{
		Movies: *movie,
	}

	// Try to get OMDB rating if the movie has an IMDB ID
	if movie.ImdbID != "" {
		omdbData, err := fetch.RatingInfo(movie.ImdbID)
		if err == nil && omdbData != nil {
			combinedResponse.OMDBRating = omdbData.ImdbRating
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(combinedResponse)
}
