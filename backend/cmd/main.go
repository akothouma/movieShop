package main

import (
	"fmt"
	handlers "movieshop/backend/Handlers"
	dependencies "movieshop/backend/cmd/web/dependancies"
	"net/http"
)

func main() {
	db, err := InitDB()
	if err != nil {
		fmt.Errorf("Failed to initialize database")
		return
	}

	cache, err := InitCache()
	if err != nil {
		fmt.Errorf("Failed to initialize cache")
		return
	}

	dep := &dependencies.Dependencies{
		DB:    db,
		Cache: cache,
	}

	mux := http.NewServeMux()
	mux.Handle("/users", handlers.UserHandler(dep))
	mux.Handle("/userPrefs", handlers.PreferencesHandler(dep))
	mux.Handle("/",handlers.MovieListHandler(dep))
	mux.Handle("/trending", handlers.TrendingListHandler(dep))
	mux.Handle("/filter", handlers.FilterListHandler(dep))
	mux.Handle("/reccomendation", handlers.RecommendationHandler(dep))
}
