package handlers

import (
	"encoding/json"
	dependencies "movieshop/backend/cmd/web/dependancies"
	"net/http"
	"strings"
	"time"
)

type PreferenceHandler struct {
	dep dependencies.Dependencies
}
type preferences struct {
	Theme string `json:"theme"`
}

func PreferenceHandler(deps dependencies.Dependencies) *UserPrefHandler {
	return &PreferenceHandler{dep: deps}
}

func (p *PreferenceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID := r.context().Value("user_id")
	switch r.Method {
	case http.MethodGet:
		p.getPreferences(w, r, userID)
	case http.MethodPut:
		p.updatePreferences(w, r, userID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

	}
}

func (p *PreferenceHandler) getPreferences(w http.ResponseWriter, r *http.Request, userID string) {
	ctx := r.Context()

	preference, err := p.dep.Redis.Get(ctx, "prefs"+userID).Result()
	if err == nil {
		json.NewEncoder(w).Encode(preference)
		return
	}
	theme, err := p.dep.DB.getPreferenceFromDB(ctx, userID)
	if err != nil {
		theme = "light"
	}

	p.dep.Redis.Set(ctx, "prefs"+userID, theme, 365*24*time.Hour)
	json.NewEncoder(w).Encode(map[string]string{"theme": theme})
}

func (u *PreferenceHandler) updatePreferences(w http.ResponseWriter, r *http.Request, userID string) {
	var prefs preferences

	if err := json.NewDecoder(r.body).Decode(&prefs); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	newTheme := SaveTheme(prefs.Theme)

	u.dep.DB.SavePreferences(r.Context, userID, newTheme)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTheme)
}

func SaveTheme(theme string) string {
	switch strings.ToLower(strings.TrimSpace(theme)) {
	case "light", "dark":
		return strings.ToLower(theme)
	default:
		return "light"
	}
}
