package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var db *sql.DB

type Preference struct {
	UserID    string    `json:"user_id"`
	Theme     string    `json:"theme"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetPreferenceFromDB(ctx context.Context, userID string) (string, error) {
	var theme string

	query := "SELECT theme FROM user_preferences WHERE user_id = ?"
	err := db.QueryRowContext(ctx, query, userID).Scan(&theme)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return "light", nil
	case err != nil:
		return "", fmt.Errorf("failed to get preferences: %w", err)
	default:
		return theme, nil
	}
}

func SavePreferences(ctx context.Context, userID, theme string) error {
	_, err := db.ExecContext(ctx,
		"INSERT INTO user_preferences (user_id, theme) VALUES (?, ?) "+
			"ON CONFLICT(user_id) DO UPDATE SET theme = ?, updated_at = CURRENT_TIMESTAMP",
		userID, theme, theme)

	return err
}
