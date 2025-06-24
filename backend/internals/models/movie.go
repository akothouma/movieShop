package models

// Movies represents movie data from TMDB API
type Movies struct {
    ID           int     `json:"id"`
    Title        string  `json:"title"`
    ReleaseDate  string  `json:"release_date"`
    PosterPath   string  `json:"poster_path"`
    Overview     string  `json:"overview"`
    VoteAverage  float64 `json:"vote_average"`
    ImdbID       string  `json:"imdb_id,omitempty"`
}
