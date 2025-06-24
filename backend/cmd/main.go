package main

import (
	middlewares "command-line-arguments/home/lakoth/movieShop/backend/internals/middlewares/auth.go"
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

	// Public routes
	mux.Handle("/",middlewares.RateLimitter(handlers.MovieListHandler(dep)))
	mux.Handle("/trending", middlewares.RateLimitter(handlers.TrendingListHandler(dep)))
	mux.Handle("/filter", middlewares.RateLimitter(handlers.FilterListHandler(dep)))

	// Protected routes (with auth + rate limiting)
	protectedRoutes := map[string]http.Handler{
		"/users":               handlers.UserHandler(dep),
		"/userPrefs":           handlers.PreferencesHandler(dep),
		"/recommendation":      handlers.RecommendationHandler(dep),
		"/addToWatchlist":      handlers.AddToWatchlist(dep),
		"/removeFromWatchlist": handlers.RemoveFromWatchlist(dep),
	}

	for path, handler := range protectedRoutes {
		mux.Handle(path, middlewares.ChainMiddlewares(handler))
	}


	http.ListenAndServe(":8000", mux)
}
